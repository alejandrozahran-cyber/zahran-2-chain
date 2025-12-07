package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

// Block Explorer Backend - Indexer + REST API
// Indexes blockchain data into PostgreSQL for fast queries

type Indexer struct {
	db          *sql.DB
	latestBlock uint64
	rpcEndpoint string
}

type Block struct {
	Number      uint64    `json:"number"`
	Hash        string    `json:"hash"`
	ParentHash  string    `json:"parent_hash"`
	Timestamp   time.Time `json:"timestamp"`
	Miner       string    `json:"miner"`
	GasUsed     uint64    `json:"gas_used"`
	GasLimit    uint64    `json:"gas_limit"`
	TxCount     int       `json:"tx_count"`
	Size        uint64    `json:"size"`
	Reward      uint64    `json:"reward"`
}

type Transaction struct {
	Hash        string    `json:"hash"`
	BlockNumber uint64    `json:"block_number"`
	FromAddr    string    `json:"from"`
	ToAddr      string    `json:"to"`
	Value       uint64    `json:"value"`
	GasPrice    uint64    `json:"gas_price"`
	GasUsed     uint64    `json:"gas_used"`
	Nonce       uint64    `json:"nonce"`
	Timestamp   time.Time `json:"timestamp"`
	Status      string    `json:"status"`
	InputData   string    `json:"input_data"`
}

type Address struct {
	Address     string    `json:"address"`
	Balance     uint64    `json:"balance"`
	TxCount     int       `json:"tx_count"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
	IsContract  bool      `json:"is_contract"`
}

func NewIndexer(dbURL, rpcEndpoint string) (*Indexer, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, err
	}

	return &Indexer{
		db:          db,
		latestBlock: 0,
		rpcEndpoint: rpcEndpoint,
	}, nil
}

func createTables(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS blocks (
		number BIGINT PRIMARY KEY,
		hash VARCHAR(66) UNIQUE NOT NULL,
		parent_hash VARCHAR(66) NOT NULL,
		timestamp TIMESTAMP NOT NULL,
		miner VARCHAR(66) NOT NULL,
		gas_used BIGINT NOT NULL,
		gas_limit BIGINT NOT NULL,
		tx_count INT NOT NULL,
		size BIGINT NOT NULL,
		reward BIGINT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_blocks_timestamp ON blocks(timestamp);
	CREATE INDEX IF NOT EXISTS idx_blocks_miner ON blocks(miner);

	CREATE TABLE IF NOT EXISTS transactions (
		hash VARCHAR(66) PRIMARY KEY,
		block_number BIGINT NOT NULL,
		from_addr VARCHAR(66) NOT NULL,
		to_addr VARCHAR(66) NOT NULL,
		value BIGINT NOT NULL,
		gas_price BIGINT NOT NULL,
		gas_used BIGINT NOT NULL,
		nonce BIGINT NOT NULL,
		timestamp TIMESTAMP NOT NULL,
		status VARCHAR(20) NOT NULL,
		input_data TEXT,
		created_at TIMESTAMP DEFAULT NOW(),
		FOREIGN KEY (block_number) REFERENCES blocks(number)
	);

	CREATE INDEX IF NOT EXISTS idx_tx_block ON transactions(block_number);
	CREATE INDEX IF NOT EXISTS idx_tx_from ON transactions(from_addr);
	CREATE INDEX IF NOT EXISTS idx_tx_to ON transactions(to_addr);
	CREATE INDEX IF NOT EXISTS idx_tx_timestamp ON transactions(timestamp);

	CREATE TABLE IF NOT EXISTS addresses (
		address VARCHAR(66) PRIMARY KEY,
		balance BIGINT NOT NULL DEFAULT 0,
		tx_count INT NOT NULL DEFAULT 0,
		first_seen TIMESTAMP NOT NULL,
		last_seen TIMESTAMP NOT NULL,
		is_contract BOOLEAN DEFAULT FALSE,
		updated_at TIMESTAMP DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_addr_balance ON addresses(balance DESC);
	CREATE INDEX IF NOT EXISTS idx_addr_tx_count ON addresses(tx_count DESC);

	CREATE TABLE IF NOT EXISTS contracts (
		address VARCHAR(66) PRIMARY KEY,
		creator VARCHAR(66) NOT NULL,
		created_at_block BIGINT NOT NULL,
		bytecode TEXT NOT NULL,
		abi TEXT,
		verified BOOLEAN DEFAULT FALSE,
		FOREIGN KEY (address) REFERENCES addresses(address)
	);
	`

	_, err := db.Exec(schema)
	return err
}

// Index new block
func (idx *Indexer) IndexBlock(block *Block) error {
	tx, err := idx.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert block
	_, err = tx.Exec(`
		INSERT INTO blocks (number, hash, parent_hash, timestamp, miner, gas_used, gas_limit, tx_count, size, reward)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (number) DO NOTHING
	`, block.Number, block.Hash, block.ParentHash, block.Timestamp, block.Miner,
		block.GasUsed, block.GasLimit, block.TxCount, block.Size, block.Reward)

	if err != nil {
		return err
	}

	// Update miner address
	idx.updateAddress(tx, block.Miner, block. Timestamp)

	if err := tx.Commit(); err != nil {
		return err
	}

	idx.latestBlock = block.Number
	log.Printf("✅ Indexed block %d (hash: %s, txs: %d)", block.Number, block.Hash[:16], block.TxCount)

	return nil
}

// Index transaction
func (idx *Indexer) IndexTransaction(tx *Transaction) error {
	dbtx, err := idx.db. Begin()
	if err != nil {
		return err
	}
	defer dbtx. Rollback()

	// Insert transaction
	_, err = dbtx.Exec(`
		INSERT INTO transactions (hash, block_number, from_addr, to_addr, value, gas_price, gas_used, nonce, timestamp, status, input_data)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (hash) DO NOTHING
	`, tx.Hash, tx.BlockNumber, tx.FromAddr, tx.ToAddr, tx.Value, tx.GasPrice, tx.GasUsed, tx.Nonce, tx. Timestamp, tx.Status, tx.InputData)

	if err != nil {
		return err
	}

	// Update addresses
	idx.updateAddress(dbtx, tx.FromAddr, tx.Timestamp)
	idx.updateAddress(dbtx, tx.ToAddr, tx. Timestamp)

	// Update balances (simplified)
	_, err = dbtx.Exec("UPDATE addresses SET balance = balance - $1 WHERE address = $2", tx.Value, tx.FromAddr)
	if err != nil {
		return err
	}

	_, err = dbtx. Exec("UPDATE addresses SET balance = balance + $1 WHERE address = $2", tx.Value, tx.ToAddr)
	if err != nil {
		return err
	}

	if err := dbtx.Commit(); err != nil {
		return err
	}

	return nil
}

func (idx *Indexer) updateAddress(tx *sql.Tx, address string, timestamp time.Time) {
	tx.Exec(`
		INSERT INTO addresses (address, tx_count, first_seen, last_seen)
		VALUES ($1, 1, $2, $2)
		ON CONFLICT (address) DO UPDATE SET
			tx_count = addresses.tx_count + 1,
			last_seen = $2,
			updated_at = NOW()
	`, address, timestamp)
}

// REST API Handlers

func (idx *Indexer) ServeAPI() {
	http.HandleFunc("/api/blocks/latest", idx.handleLatestBlocks)
	http.HandleFunc("/api/block/", idx.handleGetBlock)
	http.HandleFunc("/api/tx/", idx.handleGetTransaction)
	http.HandleFunc("/api/address/", idx.handleGetAddress)
	http.HandleFunc("/api/stats", idx.handleStats)
	http.HandleFunc("/api/search", idx.handleSearch)

	// Enable CORS
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	})

	log.Println("🌐 Explorer API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (idx *Indexer) handleLatestBlocks(w http.ResponseWriter, r *http.Request) {
	rows, err := idx.db.Query(`
		SELECT number, hash, timestamp, miner, tx_count, gas_used
		FROM blocks
		ORDER BY number DESC
		LIMIT 20
	`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	blocks := []map[string]interface{}{}
	for rows.Next() {
		var num uint64
		var hash, miner string
		var timestamp time.Time
		var txCount int
		var gasUsed uint64

		rows.Scan(&num, &hash, &timestamp, &miner, &txCount, &gasUsed)

		blocks = append(blocks, map[string]interface{}{
			"number":    num,
			"hash":      hash,
			"timestamp": timestamp,
			"miner":     miner,
			"tx_count":  txCount,
			"gas_used":  gasUsed,
		})
	}

	json.NewEncoder(w).Encode(blocks)
}

func (idx *Indexer) handleGetBlock(w http.ResponseWriter, r *http.Request) {
	// Extract block number from URL
	// Implementation similar to handleLatestBlocks
	w.Write([]byte(`{"status": "ok"}`))
}

func (idx *Indexer) handleGetTransaction(w http.ResponseWriter, r *http.Request) {
	// Get transaction by hash
	w.Write([]byte(`{"status": "ok"}`))
}

func (idx *Indexer) handleGetAddress(w http. ResponseWriter, r *http.Request) {
	// Get address details
	w.Write([]byte(`{"status": "ok"}`))
}

func (idx *Indexer) handleStats(w http.ResponseWriter, r *http.Request) {
	var totalBlocks, totalTxs, totalAddresses int64

	idx.db.QueryRow("SELECT COUNT(*) FROM blocks"). Scan(&totalBlocks)
	idx.db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&totalTxs)
	idx.db.QueryRow("SELECT COUNT(*) FROM addresses").Scan(&totalAddresses)

	stats := map[string]interface{}{
		"total_blocks":    totalBlocks,
		"total_txs":       totalTxs,
		"total_addresses": totalAddresses,
		"latest_block":    idx. latestBlock,
	}

	json.NewEncoder(w).Encode(stats)
}

func (idx *Indexer) handleSearch(w http. ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	// Search blocks, txs, addresses
	w.Write([]byte(fmt.Sprintf(`{"query": "%s"}`, query)))
}

func main() {
	dbURL := "postgres://nusa:password@localhost/nusa_explorer?sslmode=disable"
	rpcEndpoint := "http://localhost:8545"

	indexer, err := NewIndexer(dbURL, rpcEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	// Start API server
	go indexer.ServeAPI()

	// Indexing loop (simplified)
	log.Println("🔍 Explorer indexer started")
	select {}
}
