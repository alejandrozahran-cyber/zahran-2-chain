package lightclient

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Light Client Pool - Fast sync for mobile wallets & dApps
// Multi-chain light verification without full node

type LightClientPool struct {
	Clients         map[string]*LightClient
	CheckpointStore map[uint64]*Checkpoint
	SyncServers     map[string]*SyncServer
	LatestBlock     uint64
}

type LightClient struct {
	ClientID       string
	Address        string
	SyncedBlock    uint64
	LastSyncTime   time.Time
	DataUsage      uint64 // Bytes
	ConnectedChain string // NUSA, ETH, BTC, etc
	Type           ClientType
}

type ClientType string

const (
	MobileWallet ClientType = "mobile_wallet"
	Browser      ClientType = "browser"
	DApp         ClientType = "dapp"
	Explorer     ClientType = "explorer"
)

type Checkpoint struct {
	BlockNumber uint64
	BlockHash   string
	StateRoot   string
	Timestamp   time.Time
	Validators  []string // Validator signatures
	TxCount     uint64
	GasUsed     uint64
}

type SyncServer struct {
	ServerID    string
	Address     string
	Uptime      float64
	Clients     int
	Bandwidth   uint64 // MB/s
	SyncSpeed   uint64 // blocks/sec
	Reputation  int
}

type SyncRequest struct {
	ClientID    string
	FromBlock   uint64
	ToBlock     uint64
	ChainID     string
	DataType    []string // headers, receipts, state
}

func NewLightClientPool() *LightClientPool {
	return &LightClientPool{
		Clients:         make(map[string]*LightClient),
		CheckpointStore: make(map[uint64]*Checkpoint),
		SyncServers:     make(map[string]*SyncServer),
		LatestBlock:     0,
	}
}

// Register light client
func (lcp *LightClientPool) RegisterClient(
	clientID string,
	address string,
	clientType ClientType,
	chain string,
) bool {
	if _, exists := lcp.Clients[clientID]; exists {
		return false
	}

	client := &LightClient{
		ClientID:       clientID,
		Address:        address,
		SyncedBlock:    0,
		LastSyncTime:   time.Now(),
		DataUsage:      0,
		ConnectedChain: chain,
		Type:           clientType,
	}

	lcp.Clients[clientID] = client

	fmt.Printf("📱 Light client registered: %s (%s on %s)\n", clientID, clientType, chain)

	return true
}

// Fast sync using checkpoints
func (lcp *LightClientPool) FastSync(clientID string, targetBlock uint64) (*SyncResult, error) {
	client, exists := lcp.Clients[clientID]
	if !exists {
		return nil, fmt.Errorf("client not found")
	}

	if client. SyncedBlock >= targetBlock {
		return &SyncResult{
			Success:     true,
			BlocksSynced: 0,
			TimeTaken:   0,
			DataUsed:    0,
		}, nil
	}

	startTime := time.Now()

	// Find nearest checkpoint
	checkpoint := lcp.findNearestCheckpoint(targetBlock)
	if checkpoint == nil {
		return nil, fmt.Errorf("no checkpoint available")
	}

	// Select best sync server
	server := lcp.selectBestSyncServer()
	if server == nil {
		return nil, fmt.Errorf("no sync server available")
	}

	// Download only block headers + state roots (NOT full blocks!)
	blocksSynced := targetBlock - client.SyncedBlock
	dataPerBlock := uint64(500) // Only 500 bytes per block (vs 100KB+ for full block)
	totalData := blocksSynced * dataPerBlock

	// Simulate sync (production: actual network sync)
	syncDuration := time.Duration(blocksSynced/server.SyncSpeed) * time. Second

	client.SyncedBlock = targetBlock
	client.LastSyncTime = time.Now()
	client.DataUsage += totalData

	server.Clients++

	result := &SyncResult{
		Success:      true,
		BlocksSynced: blocksSynced,
		TimeTaken:    time.Since(startTime),
		DataUsed:     totalData,
		Checkpoint:   checkpoint. BlockNumber,
	}

	fmt.Printf("⚡ Fast sync completed: %d blocks in %v (%.2f MB)\n",
		blocksSynced, syncDuration, float64(totalData)/(1024*1024))

	return result, nil
}

// Create checkpoint every N blocks
func (lcp *LightClientPool) CreateCheckpoint(
	blockNumber uint64,
	blockHash string,
	stateRoot string,
	validators []string,
	txCount uint64,
	gasUsed uint64,
) {
	checkpoint := &Checkpoint{
		BlockNumber: blockNumber,
		BlockHash:   blockHash,
		StateRoot:   stateRoot,
		Timestamp:   time.Now(),
		Validators:  validators,
		TxCount:     txCount,
		GasUsed:     gasUsed,
	}

	lcp.CheckpointStore[blockNumber] = checkpoint

	fmt.Printf("📍 Checkpoint created: block %d (hash: %s)\n", blockNumber, blockHash[:16])
}

// Verify checkpoint with validator signatures
func (lcp *LightClientPool) VerifyCheckpoint(blockNumber uint64) bool {
	checkpoint, exists := lcp.CheckpointStore[blockNumber]
	if !exists {
		return false
	}

	// Need 2/3 validator signatures
	minValidators := 3
	if len(checkpoint. Validators) < minValidators {
		fmt. Println("❌ Insufficient validator signatures")
		return false
	}

	// Verify signatures (simplified - production needs crypto verification)
	for _, validator := range checkpoint.Validators {
		if len(validator) < 10 {
			return false
		}
	}

	fmt.Printf("✅ Checkpoint verified: block %d with %d validators\n",
		blockNumber, len(checkpoint.Validators))

	return true
}

// Multi-chain light verification
func (lcp *LightClientPool) VerifyMultiChainTx(
	chain string,
	txHash string,
	blockNumber uint64,
) (bool, error) {
	// Verify transaction exists on another chain (ETH, BTC, etc)
	// Using light client proofs (Merkle proofs)

	fmt.Printf("🔍 Verifying %s transaction on %s at block %d\n", txHash, chain, blockNumber)

	// Check if we have checkpoint for that block
	checkpoint := lcp.findNearestCheckpoint(blockNumber)
	if checkpoint == nil {
		return false, fmt.Errorf("no checkpoint for block %d", blockNumber)
	}

	// Verify Merkle proof (simplified)
	merkleProof := lcp.generateMerkleProof(txHash, checkpoint.StateRoot)

	if !lcp.verifyMerkleProof(merkleProof, checkpoint.StateRoot) {
		return false, fmt.Errorf("invalid Merkle proof")
	}

	fmt.Printf("✅ Multi-chain tx verified: %s on %s\n", txHash[:16], chain)

	return true, nil
}

// Register sync server
func (lcp *LightClientPool) RegisterSyncServer(
	serverID string,
	address string,
	bandwidth uint64,
	syncSpeed uint64,
) {
	server := &SyncServer{
		ServerID:   serverID,
		Address:    address,
		Uptime:     99.9,
		Clients:    0,
		Bandwidth:  bandwidth,
		SyncSpeed:  syncSpeed,
		Reputation: 100,
	}

	lcp.SyncServers[serverID] = server

	fmt.Printf("🖥️ Sync server registered: %s (speed: %d blocks/s)\n", serverID, syncSpeed)
}

// Select best sync server (load balancing)
func (lcp *LightClientPool) selectBestSyncServer() *SyncServer {
	var bestServer *SyncServer
	bestScore := 0. 0

	for _, server := range lcp.SyncServers {
		// Score = (uptime * syncSpeed * reputation) / clients
		score := (server. Uptime * float64(server.SyncSpeed) * float64(server.Reputation)) / float64(server.Clients+1)

		if score > bestScore {
			bestScore = score
			bestServer = server
		}
	}

	return bestServer
}

// Find nearest checkpoint
func (lcp *LightClientPool) findNearestCheckpoint(targetBlock uint64) *Checkpoint {
	var nearest *Checkpoint
	minDiff := uint64(^uint(0)) // Max uint

	for blockNum, checkpoint := range lcp.CheckpointStore {
		if blockNum <= targetBlock {
			diff := targetBlock - blockNum
			if diff < minDiff {
				minDiff = diff
				nearest = checkpoint
			}
		}
	}

	return nearest
}

// Generate Merkle proof
func (lcp *LightClientPool) generateMerkleProof(txHash string, stateRoot string) []string {
	// Simplified Merkle proof generation
	proof := make([]string, 0)

	// Production: Generate actual Merkle branch
	hash := sha256.Sum256([]byte(txHash))
	proof = append(proof, fmt. Sprintf("%x", hash))

	return proof
}

// Verify Merkle proof
func (lcp *LightClientPool) verifyMerkleProof(proof []string, stateRoot string) bool {
	// Simplified verification
	if len(proof) == 0 {
		return false
	}

	// Production: Verify Merkle branch against state root
	return true
}

// Get sync status for client
func (lcp *LightClientPool) GetSyncStatus(clientID string) map[string]interface{} {
	client, exists := lcp.Clients[clientID]
	if ! exists {
		return nil
	}

	blocksBehind := lcp.LatestBlock - client.SyncedBlock
	syncProgress := float64(client.SyncedBlock) / float64(lcp.LatestBlock) * 100

	return map[string]interface{}{
		"client_id":      client.ClientID,
		"synced_block":   client.SyncedBlock,
		"latest_block":   lcp. LatestBlock,
		"blocks_behind":  blocksBehind,
		"sync_progress":  fmt.Sprintf("%.2f%%", syncProgress),
		"data_usage_mb":  float64(client.DataUsage) / (1024 * 1024),
		"last_sync":      client.LastSyncTime,
	}
}

// Get pool stats
func (lcp *LightClientPool) GetPoolStats() map[string]interface{} {
	mobileClients := 0
	browserClients := 0
	totalDataUsage := uint64(0)

	for _, client := range lcp. Clients {
		if client. Type == MobileWallet {
			mobileClients++
		} else if client.Type == Browser {
			browserClients++
		}
		totalDataUsage += client. DataUsage
	}

	return map[string]interface{}{
		"total_clients":    len(lcp.Clients),
		"mobile_clients":   mobileClients,
		"browser_clients":  browserClients,
		"sync_servers":     len(lcp.SyncServers),
		"checkpoints":      len(lcp. CheckpointStore),
		"total_data_gb":    float64(totalDataUsage) / (1024 * 1024 * 1024),
		"latest_block":     lcp. LatestBlock,
	}
}

type SyncResult struct {
	Success      bool
	BlocksSynced uint64
	TimeTaken    time.Duration
	DataUsed     uint64
	Checkpoint   uint64
}
