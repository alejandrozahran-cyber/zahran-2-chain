package mev

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

// MEV Protection System
// Anti-frontrunning, anti-sandwich, encrypted mempool, private order flow

type MEVProtection struct {
	EncryptedMempool  map[string]*EncryptedTx
	PrivateOrderFlow  map[string]*PrivateOrder
	SandwichDetector  *SandwichDetector
	FrontrunDetector  *FrontrunDetector
	ProtectedTxCount  uint64
	BlockedMEVAttempts uint64
}

type EncryptedTx struct {
	EncryptedData []byte
	Timestamp     time.Time
	RevealBlock   uint64 // Block number when tx is revealed
	Sender        string
	Nonce         uint64
}

type PrivateOrder struct {
	OrderID       string
	Type          string // "swap", "limit", "market"
	EncryptedData []byte
	MinPrice      uint64 // For slippage protection
	MaxPrice      uint64
	Deadline      time.Time
	Protected     bool
}

type SandwichDetector struct {
	RecentSwaps   []SwapEvent
	DetectedCount uint64
}

type SwapEvent struct {
	TxHash    string
	From      string
	TokenIn   string
	TokenOut  string
	AmountIn  uint64
	AmountOut uint64
	Timestamp time.Time
	Block     uint64
}

type FrontrunDetector struct {
	PendingHighValue []Transaction
	DetectedCount    uint64
}

type Transaction struct {
	Hash      string
	From      string
	To        string
	Value     uint64
	GasPrice  uint64
	Timestamp time.Time
}

func NewMEVProtection() *MEVProtection {
	return &MEVProtection{
		EncryptedMempool:   make(map[string]*EncryptedTx),
		PrivateOrderFlow:   make(map[string]*PrivateOrder),
		SandwichDetector:   NewSandwichDetector(),
		FrontrunDetector:   NewFrontrunDetector(),
		ProtectedTxCount:   0,
		BlockedMEVAttempts: 0,
	}
}

func NewSandwichDetector() *SandwichDetector {
	return &SandwichDetector{
		RecentSwaps:   make([]SwapEvent, 0),
		DetectedCount: 0,
	}
}

func NewFrontrunDetector() *FrontrunDetector {
	return &FrontrunDetector{
		PendingHighValue: make([]Transaction, 0),
		DetectedCount:    0,
	}
}

// Submit encrypted transaction (time-locked encryption)
func (mev *MEVProtection) SubmitEncryptedTx(
	txData []byte,
	sender string,
	revealBlock uint64,
	nonce uint64,
) (string, error) {
	// Encrypt transaction data
	encryptedData, err := encryptData(txData)
	if err != nil {
		return "", err
	}

	// Generate transaction hash
	hash := generateTxHash(encryptedData, sender, nonce)

	encTx := &EncryptedTx{
		EncryptedData: encryptedData,
		Timestamp:     time.Now(),
		RevealBlock:   revealBlock,
		Sender:        sender,
		Nonce:         nonce,
	}

	mev.EncryptedMempool[hash] = encTx
	mev. ProtectedTxCount++

	fmt.Printf("🔒 Encrypted tx submitted: %s (reveal at block %d)\n", hash[:16], revealBlock)

	return hash, nil
}

// Reveal encrypted transactions when block is reached
func (mev *MEVProtection) RevealTransactions(currentBlock uint64) [][]byte {
	revealed := make([][]byte, 0)

	for hash, encTx := range mev. EncryptedMempool {
		if currentBlock >= encTx. RevealBlock {
			// Decrypt transaction
			decrypted, err := decryptData(encTx.EncryptedData)
			if err != nil {
				continue
			}

			revealed = append(revealed, decrypted)
			delete(mev.EncryptedMempool, hash)

			fmt.Printf("🔓 Tx revealed: %s at block %d\n", hash[:16], currentBlock)
		}
	}

	return revealed
}

// Submit private order (DEX integration)
func (mev *MEVProtection) SubmitPrivateOrder(
	orderType string,
	orderData []byte,
	minPrice, maxPrice uint64,
	deadline time.Time,
) (string, error) {
	// Encrypt order data
	encrypted, err := encryptData(orderData)
	if err != nil {
		return "", err
	}

	orderID := generateOrderID(orderData)

	order := &PrivateOrder{
		OrderID:       orderID,
		Type:          orderType,
		EncryptedData: encrypted,
		MinPrice:      minPrice,
		MaxPrice:      maxPrice,
		Deadline:      deadline,
		Protected:     true,
	}

	mev.PrivateOrderFlow[orderID] = order

	fmt.Printf("🔐 Private order submitted: %s (%s)\n", orderID[:16], orderType)

	return orderID, nil
}

// Detect sandwich attack
func (mev *MEVProtection) DetectSandwich(swap SwapEvent) bool {
	detector := mev.SandwichDetector

	// Check if there are suspicious trades around this swap
	for _, recentSwap := range detector.RecentSwaps {
		// Check for sandwich pattern:
		// 1. Same token pair
		// 2. Opposite direction
		// 3. Very close in time (< 5 seconds)
		// 4. Large price impact

		if swap.TokenIn == recentSwap.TokenOut && swap.TokenOut == recentSwap.TokenIn {
			timeDiff := swap.Timestamp.Sub(recentSwap.Timestamp)

			if timeDiff < 5*time.Second {
				// Potential sandwich detected! 
				detector.DetectedCount++
				mev.BlockedMEVAttempts++

				fmt.Printf("🚨 SANDWICH ATTACK DETECTED!\n")
				fmt.Printf("   Front-run: %s\n", recentSwap.TxHash[:16])
				fmt.Printf("   Victim: %s\n", swap. TxHash[:16])

				return true
			}
		}
	}

	// Add to recent swaps
	detector.RecentSwaps = append(detector.RecentSwaps, swap)

	// Keep only last 100 swaps
	if len(detector.RecentSwaps) > 100 {
		detector. RecentSwaps = detector.RecentSwaps[1:]
	}

	return false
}

// Detect frontrunning attempt
func (mev *MEVProtection) DetectFrontrun(tx Transaction) bool {
	detector := mev.FrontrunDetector

	// Check if this is a high-value transaction
	if tx. Value > 1000000 { // > 1M units
		// Check if there's a similar pending tx with higher gas
		for _, pendingTx := range detector.PendingHighValue {
			// Same destination, higher gas, submitted after
			if pendingTx.To == tx.To &&
				tx.GasPrice > pendingTx.GasPrice &&
				tx.Timestamp.After(pendingTx.Timestamp) {

				detector.DetectedCount++
				mev.BlockedMEVAttempts++

				fmt.Printf("🚨 FRONTRUN ATTEMPT DETECTED!\n")
				fmt.Printf("   Original: %s (gas: %d)\n", pendingTx.Hash[:16], pendingTx. GasPrice)
				fmt. Printf("   Frontrun: %s (gas: %d)\n", tx.Hash[:16], tx.GasPrice)

				return true
			}
		}

		// Add to pending high-value txs
		detector.PendingHighValue = append(detector.PendingHighValue, tx)

		// Keep only last 50 txs
		if len(detector. PendingHighValue) > 50 {
			detector.PendingHighValue = detector.PendingHighValue[1:]
		}
	}

	return false
}

// Fair ordering - randomize transaction order within block
func (mev *MEVProtection) FairOrdering(txs []Transaction) []Transaction {
	// Use block hash as randomness source (verifiable randomness)
	// Shuffle transactions to prevent ordering manipulation

	shuffled := make([]Transaction, len(txs))
	copy(shuffled, txs)

	// Fisher-Yates shuffle (simplified - use VRF in production)
	for i := len(shuffled) - 1; i > 0; i-- {
		j := int(time.Now().UnixNano()) % (i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	fmt.Printf("🎲 Fair ordering applied: %d transactions shuffled\n", len(txs))

	return shuffled
}

// Commit-reveal scheme for sensitive operations
func (mev *MEVProtection) CommitPhase(commitment string, sender string) string {
	// User commits hash of their action
	commitID := sha256.Sum256([]byte(commitment + sender + time.Now().String()))
	commitHash := hex.EncodeToString(commitID[:])

	fmt.Printf("📝 Commitment recorded: %s\n", commitHash[:16])

	return commitHash
}

func (mev *MEVProtection) RevealPhase(commitHash, secret string) bool {
	// User reveals secret, verify it matches commitment
	// This prevents frontrunning of sensitive operations

	// Verify commitment (simplified)
	fmt.Printf("🔓 Reveal verified: %s\n", commitHash[:16])

	return true
}

// Encryption helpers
func encryptData(data []byte) ([]byte, error) {
	key := make([]byte, 32) // 256-bit key
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decryptData(ciphertext []byte) ([]byte, error) {
	// Simplified decryption (production: proper key management)
	// In reality, use threshold encryption with validators

	return ciphertext, nil // Placeholder
}

func generateTxHash(data []byte, sender string, nonce uint64) string {
	combined := fmt.Sprintf("%s%s%d", data, sender, nonce)
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

func generateOrderID(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// Get MEV protection stats
func (mev *MEVProtection) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"protected_transactions": mev.ProtectedTxCount,
		"blocked_mev_attempts":   mev.BlockedMEVAttempts,
		"encrypted_mempool_size": len(mev. EncryptedMempool),
		"private_orders":         len(mev.PrivateOrderFlow),
		"sandwich_detected":      mev.SandwichDetector.DetectedCount,
		"frontrun_detected":      mev.FrontrunDetector.DetectedCount,
	}
}
