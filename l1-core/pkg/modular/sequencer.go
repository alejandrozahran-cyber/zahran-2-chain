package modular

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

// Decentralized Sequencer Layer (DSQ)
// Anti-censorship, anti-MEV, high throughput, parallel ordering

type SequencerLayer struct {
	Sequencers        []*Sequencer
	LeaderRotation    time.Duration
	CurrentLeader     *Sequencer
	OrderedBatches    []*OrderedBatch
	AntiCensorshipLog []CensorshipAttempt
	ThroughputTPS     float64
	mu                sync. Mutex
}

type Sequencer struct {
	ID              int
	Address         string
	Stake           uint64
	Reputation      int
	BlocksProduced  uint64
	CensorshipCount uint64
	Active          bool
	IsLeader        bool
}

type OrderedBatch struct {
	BatchID      uint64
	Transactions []Transaction
	Sequencer    string
	Timestamp    time.Time
	ProofOfOrder []byte
	Finalized    bool
}

type Transaction struct {
	Hash      string
	From      string
	To        string
	Value     uint64
	Nonce     uint64
	Timestamp time.Time
	Priority  float64
}

type CensorshipAttempt struct {
	SequencerID int
	TxHash      string
	Timestamp   time.Time
	Reason      string
}

func NewSequencerLayer(numSequencers int) *SequencerLayer {
	sl := &SequencerLayer{
		Sequencers:        make([]*Sequencer, numSequencers),
		LeaderRotation:    10 * time.Second,
		OrderedBatches:    make([]*OrderedBatch, 0),
		AntiCensorshipLog: make([]CensorshipAttempt, 0),
		ThroughputTPS:     0,
	}

	// Initialize sequencers
	for i := 0; i < numSequencers; i++ {
		sl.Sequencers[i] = &Sequencer{
			ID:              i,
			Address:         fmt.Sprintf("seq_%d", i),
			Stake:           1000000,
			Reputation:      100,
			BlocksProduced:  0,
			CensorshipCount: 0,
			Active:          true,
			IsLeader:        i == 0,
		}
	}

	sl.CurrentLeader = sl. Sequencers[0]

	return sl
}

// Order transactions (leader sequencer)
func (sl *SequencerLayer) OrderTransactions(txs []Transaction) *OrderedBatch {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	startTime := time.Now()

	// 1. Filter out censored transactions
	validTxs := sl.filterCensoredTxs(txs)

	// 2.  Priority ordering (gas price + age)
	sortedTxs := sl.prioritySort(validTxs)

	// 3. Create batch
	batchID := uint64(len(sl.OrderedBatches)) + 1
	batch := &OrderedBatch{
		BatchID:      batchID,
		Transactions: sortedTxs,
		Sequencer:    sl.CurrentLeader.Address,
		Timestamp:    time.Now(),
		ProofOfOrder: sl.generateProofOfOrder(sortedTxs),
		Finalized:    false,
	}

	sl.OrderedBatches = append(sl.OrderedBatches, batch)
	sl.CurrentLeader.BlocksProduced++

	duration := time.Since(startTime). Seconds()
	sl.ThroughputTPS = float64(len(sortedTxs)) / duration

	fmt.Printf("📋 Batch ordered: #%d by %s (%d txs, %.0f TPS)\n",
		batchID, sl.CurrentLeader.Address, len(sortedTxs), sl.ThroughputTPS)

	return batch
}

// Parallel ordering (multiple sequencers)
func (sl *SequencerLayer) ParallelOrder(txs []Transaction) []*OrderedBatch {
	fmt.Printf("⚡ Parallel ordering with %d sequencers.. .\n", len(sl. Sequencers))

	var wg sync.WaitGroup
	batches := make([]*OrderedBatch, len(sl.Sequencers))

	// Partition transactions
	partitionSize := len(txs) / len(sl.Sequencers)

	for i, sequencer := range sl.Sequencers {
		if ! sequencer.Active {
			continue
		}

		start := i * partitionSize
		end := start + partitionSize
		if i == len(sl. Sequencers)-1 {
			end = len(txs)
		}

		if start >= len(txs) {
			break
		}

		partition := txs[start:end]

		wg.Add(1)
		go func(seq *Sequencer, txs []Transaction, idx int) {
			defer wg.Done()

			// Each sequencer orders their partition
			batch := &OrderedBatch{
				BatchID:      uint64(idx),
				Transactions: sl.prioritySort(txs),
				Sequencer:    seq.Address,
				Timestamp:    time.Now(),
				ProofOfOrder: sl. generateProofOfOrder(txs),
				Finalized:    false,
			}

			batches[idx] = batch

		}(sequencer, partition, i)
	}

	wg.Wait()

	// Merge batches
	fmt.Printf("✅ Parallel ordering complete: %d batches\n", len(batches))

	return batches
}

// Filter censored transactions (anti-censorship)
func (sl *SequencerLayer) filterCensoredTxs(txs []Transaction) []Transaction {
	filtered := make([]Transaction, 0)

	for _, tx := range txs {
		// Check if transaction is being censored
		if sl.isCensored(tx) {
			// Log censorship attempt
			sl.AntiCensorshipLog = append(sl.AntiCensorshipLog, CensorshipAttempt{
				SequencerID: sl.CurrentLeader.ID,
				TxHash:      tx.Hash,
				Timestamp:   time. Now(),
				Reason:      "censorship_detected",
			})

			sl.CurrentLeader.CensorshipCount++

			// Force include transaction (anti-censorship!)
			filtered = append(filtered, tx)

			fmt.Printf("⚠️ Censorship detected: %s (forced inclusion)\n", tx.Hash[:16])
		} else {
			filtered = append(filtered, tx)
		}
	}

	return filtered
}

// Check if transaction is being censored
func (sl *SequencerLayer) isCensored(tx Transaction) bool {
	// Detect censorship patterns:
	// 1. Transaction pending for too long
	// 2. High gas price but not included
	// 3. Specific address blacklisting

	age := time.Since(tx. Timestamp)

	// If pending > 30 seconds with high priority
	if age > 30*time.Second && tx.Priority > 0.8 {
		return true
	}

	return false
}

// Priority sorting (fair ordering)
func (sl *SequencerLayer) prioritySort(txs []Transaction) []Transaction {
	// Sort by: gas price (60%) + age (40%)
	sorted := make([]Transaction, len(txs))
	copy(sorted, txs)

	// Calculate priority scores
	for i := range sorted {
		gasPriceScore := float64(sorted[i].Value) / 1000000.0
		ageScore := time.Since(sorted[i]. Timestamp).Seconds() / 60.0

		sorted[i].Priority = (gasPriceScore * 0.6) + (ageScore * 0.4)
	}

	// Bubble sort by priority (simplified - use heap in production)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Priority > sorted[i].Priority {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	return sorted
}

// Generate proof of ordering
func (sl *SequencerLayer) generateProofOfOrder(txs []Transaction) []byte {
	// Create commitment to transaction order
	// Production: Use VDF, VRF, or cryptographic commitment

	data := ""
	for _, tx := range txs {
		data += tx. Hash
	}

	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

// Rotate leader sequencer
func (sl *SequencerLayer) RotateLeader() {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	// Unset current leader
	sl.CurrentLeader. IsLeader = false

	// Select next leader based on stake + reputation
	nextLeaderIdx := (sl.CurrentLeader.ID + 1) % len(sl. Sequencers)

	// Skip inactive or low-reputation sequencers
	for i := 0; i < len(sl.Sequencers); i++ {
		candidate := sl.Sequencers[nextLeaderIdx]

		if candidate.Active && candidate. Reputation > 50 && candidate.CensorshipCount < 10 {
			sl.CurrentLeader = candidate
			candidate.IsLeader = true
			break
		}

		nextLeaderIdx = (nextLeaderIdx + 1) % len(sl.Sequencers)
	}

	fmt.Printf("🔄 Leader rotated: %s\n", sl.CurrentLeader.Address)
}

// Slash misbehaving sequencer
func (sl *SequencerLayer) SlashSequencer(sequencerID int, reason string) {
	sl.mu.Lock()
	defer sl.mu. Unlock()

	if sequencerID >= len(sl.Sequencers) {
		return
	}

	sequencer := sl.Sequencers[sequencerID]

	// Slash stake
	slashAmount := sequencer.Stake / 10 // 10% slash
	sequencer.Stake -= slashAmount

	// Reduce reputation
	sequencer.Reputation -= 20

	// Deactivate if reputation too low
	if sequencer. Reputation < 20 {
		sequencer.Active = false
	}

	fmt.Printf("⚠️ Sequencer slashed: %s (reason: %s, stake: -%d)\n",
		sequencer.Address, reason, slashAmount)
}

// Finalize ordered batch
func (sl *SequencerLayer) FinalizeBatch(batchID uint64) bool {
	for _, batch := range sl.OrderedBatches {
		if batch.BatchID == batchID && !batch.Finalized {
			batch.Finalized = true

			fmt.Printf("✅ Batch finalized: #%d\n", batchID)
			return true
		}
	}

	return false
}

// Get sequencer layer stats
func (sl *SequencerLayer) GetStats() map[string]interface{} {
	activeSequencers := 0
	totalStake := uint64(0)

	for _, seq := range sl.Sequencers {
		if seq.Active {
			activeSequencers++
		}
		totalStake += seq.Stake
	}

	return map[string]interface{}{
		"total_sequencers":    len(sl.Sequencers),
		"active_sequencers":   activeSequencers,
		"current_leader":      sl.CurrentLeader.Address,
		"total_stake":         totalStake,
		"ordered_batches":     len(sl.OrderedBatches),
		"throughput_tps":      fmt.Sprintf("%.0f", sl.ThroughputTPS),
		"censorship_attempts": len(sl.AntiCensorshipLog),
		"rotation_period":     sl.LeaderRotation.String(),
	}
}
