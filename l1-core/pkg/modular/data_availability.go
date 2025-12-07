package modular

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

// Enhanced Data Availability Layer (Celestia-style)
// Features: Data sampling, erasure coding, KZG commitments

type EnhancedDALayer struct {
	Blocks           map[uint64]*DABlock
	Validators       []DAValidator
	SamplingNodes    []SamplingNode
	ErasureCoding    *ErasureCodingScheme
	KZGCommitments   map[uint64][]byte
	TotalDataPosted  uint64
	TotalSampled     uint64
}

type DABlock struct {
	Height       uint64
	Data         []byte
	DataShards   [][]byte  // Erasure coded shards
	ParityShards [][]byte
	Commitment   []byte    // KZG commitment
	Timestamp    time.Time
	Available    bool
}

type DAValidator struct {
	Address      string
	Stake        uint64
	StoredShards map[uint64][]int  // block -> shard indices
	Uptime       float64
}

type SamplingNode struct {
	ID           string
	SampledCount uint64
	SuccessRate  float64
}

type ErasureCodingScheme struct {
	DataShards   int  // k shards
	ParityShards int  // m shards
	TotalShards  int  // k + m
	Redundancy   float64
}

func NewEnhancedDALayer() *EnhancedDALayer {
	return &EnhancedDALayer{
		Blocks:          make(map[uint64]*DABlock),
		Validators:      make([]DAValidator, 0),
		SamplingNodes:   make([]SamplingNode, 0),
		ErasureCoding:   NewErasureCoding(128, 128),  // 50% redundancy
		KZGCommitments:  make(map[uint64][]byte),
		TotalDataPosted: 0,
		TotalSampled:    0,
	}
}

func NewErasureCoding(dataShards, parityShards int) *ErasureCodingScheme {
	return &ErasureCodingScheme{
		DataShards:   dataShards,
		ParityShards: parityShards,
		TotalShards:  dataShards + parityShards,
		Redundancy:   float64(parityShards) / float64(dataShards),
	}
}

// Post data to DA layer
func (da *EnhancedDALayer) PostData(height uint64, data []byte) (*DABlock, error) {
	fmt.Printf("💾 Posting data: block %d (%d bytes)\n", height, len(data))

	// 1. Erasure encode data
	dataShards, parityShards := da.ErasureCoding. Encode(data)

	// 2. Generate KZG commitment
	commitment := da.generateKZGCommitment(data)

	// 3. Create DA block
	block := &DABlock{
		Height:       height,
		Data:         data,
		DataShards:   dataShards,
		ParityShards: parityShards,
		Commitment:   commitment,
		Timestamp:    time.Now(),
		Available:    true,
	}

	da.Blocks[height] = block
	da.KZGCommitments[height] = commitment
	da.TotalDataPosted += uint64(len(data))

	// 4.  Distribute shards to validators
	da.distributeShards(block)

	fmt.Printf("✅ Data posted: block %d (%d data shards + %d parity shards)\n",
		height, len(dataShards), len(parityShards))

	return block, nil
}

// Data Availability Sampling (DAS)
func (da *EnhancedDALayer) SampleData(height uint64, sampleCount int) (bool, error) {
	block, exists := da.Blocks[height]
	if !exists {
		return false, fmt. Errorf("block not found")
	}

	fmt.Printf("🔍 Sampling block %d (%d random samples)...\n", height, sampleCount)

	totalShards := len(block.DataShards) + len(block.ParityShards)
	successfulSamples := 0

	// Random sampling
	for i := 0; i < sampleCount; i++ {
		shardIndex := rand.Intn(totalShards)

		// Try to retrieve shard
		if da.retrieveShard(height, shardIndex) {
			successfulSamples++
		}

		da.TotalSampled++
	}

	// Confidence threshold: 95% of samples must succeed
	confidence := float64(successfulSamples) / float64(sampleCount)
	available := confidence >= 0.95

	if available {
		fmt.Printf("✅ Block available: %d samples (%.1f%% confidence)\n",
			successfulSamples, confidence*100)
	} else {
		fmt.Printf("❌ Block NOT available: only %d/%d samples found\n",
			successfulSamples, sampleCount)
	}

	return available, nil
}

// Light client sampling (minimal bandwidth)
func (da *EnhancedDALayer) LightClientSample(height uint64) bool {
	// Light clients only sample 20 random shards
	available, _ := da.SampleData(height, 20)
	return available
}

// Reconstruct data from shards (if some missing)
func (da *EnhancedDALayer) ReconstructData(height uint64) ([]byte, error) {
	block, exists := da.Blocks[height]
	if !exists {
		return nil, fmt. Errorf("block not found")
	}

	// Check if we have enough shards (need at least k out of k+m)
	availableShards := len(block. DataShards)
	requiredShards := da.ErasureCoding.DataShards

	if availableShards < requiredShards {
		return nil, fmt.Errorf("insufficient shards: %d/%d", availableShards, requiredShards)
	}

	// Reconstruct original data
	reconstructed := da.ErasureCoding.Decode(block.DataShards, block. ParityShards)

	fmt.Printf("♻️ Data reconstructed: block %d (%d bytes)\n", height, len(reconstructed))

	return reconstructed, nil
}

// Distribute shards to validators
func (da *EnhancedDALayer) distributeShards(block *DABlock) {
	totalShards := len(block.DataShards) + len(block. ParityShards)
	shardsPerValidator := totalShards / len(da.Validators)

	if shardsPerValidator == 0 {
		shardsPerValidator = 1
	}

	for i, validator := range da.Validators {
		if validator.StoredShards == nil {
			da.Validators[i].StoredShards = make(map[uint64][]int)
		}

		// Assign shards to this validator
		shardIndices := make([]int, 0)
		for j := 0; j < shardsPerValidator; j++ {
			shardIdx := (i * shardsPerValidator) + j
			if shardIdx < totalShards {
				shardIndices = append(shardIndices, shardIdx)
			}
		}

		da.Validators[i].StoredShards[block.Height] = shardIndices

		fmt.Printf("  📤 Validator %s storing %d shards\n", validator.Address, len(shardIndices))
	}
}

// Retrieve shard from validators
func (da *EnhancedDALayer) retrieveShard(height uint64, shardIndex int) bool {
	// Find validators storing this shard
	for _, validator := range da.Validators {
		if shards, exists := validator.StoredShards[height]; exists {
			for _, idx := range shards {
				if idx == shardIndex {
					// Shard found! 
					return true
				}
			}
		}
	}

	return false
}

// Generate KZG commitment (polynomial commitment)
func (da *EnhancedDALayer) generateKZGCommitment(data []byte) []byte {
	// Simplified KZG commitment
	// Production: Use proper polynomial commitment scheme

	hash := sha256.Sum256(data)
	return hash[:]
}

// Verify KZG commitment
func (da *EnhancedDALayer) VerifyCommitment(height uint64, data []byte) bool {
	expectedCommitment := da.KZGCommitments[height]
	if expectedCommitment == nil {
		return false
	}

	actualCommitment := da.generateKZGCommitment(data)

	// Compare commitments
	if len(expectedCommitment) != len(actualCommitment) {
		return false
	}

	for i := range expectedCommitment {
		if expectedCommitment[i] != actualCommitment[i] {
			return false
		}
	}

	return true
}

// Erasure coding implementation
func (ec *ErasureCodingScheme) Encode(data []byte) ([][]byte, [][]byte) {
	// Simplified erasure coding
	// Production: Use proper Reed-Solomon encoding

	shardSize := (len(data) + ec.DataShards - 1) / ec.DataShards
	dataShards := make([][]byte, ec.DataShards)

	// Split data into k shards
	for i := 0; i < ec.DataShards; i++ {
		start := i * shardSize
		end := start + shardSize
		if end > len(data) {
			end = len(data)
		}

		shard := make([]byte, shardSize)
		copy(shard, data[start:end])
		dataShards[i] = shard
	}

	// Generate m parity shards (simplified)
	parityShards := make([][]byte, ec. ParityShards)
	for i := 0; i < ec.ParityShards; i++ {
		parity := make([]byte, shardSize)
		// XOR of all data shards (simplified - use proper RS encoding)
		for j := 0; j < shardSize; j++ {
			for k := 0; k < ec.DataShards && k < len(dataShards); k++ {
				if j < len(dataShards[k]) {
					parity[j] ^= dataShards[k][j]
				}
			}
		}
		parityShards[i] = parity
	}

	return dataShards, parityShards
}

func (ec *ErasureCodingScheme) Decode(dataShards, parityShards [][]byte) []byte {
	// Reconstruct original data from shards
	// Production: Use proper Reed-Solomon decoding

	reconstructed := make([]byte, 0)
	for _, shard := range dataShards {
		reconstructed = append(reconstructed, shard...)
	}

	return reconstructed
}

// Register DA validator
func (da *EnhancedDALayer) RegisterValidator(address string, stake uint64) {
	validator := DAValidator{
		Address:      address,
		Stake:        stake,
		StoredShards: make(map[uint64][]int),
		Uptime:       99.9,
	}

	da.Validators = append(da.Validators, validator)

	fmt.Printf("🖥️ DA Validator registered: %s (stake: %d)\n", address, stake)
}

// Get DA layer stats
func (da *EnhancedDALayer) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"total_blocks":      len(da.Blocks),
		"total_data_posted": da.TotalDataPosted,
		"total_sampled":     da.TotalSampled,
		"validators":        len(da.Validators),
		"sampling_nodes":    len(da. SamplingNodes),
		"erasure_coding": map[string]interface{}{
			"data_shards":   da.ErasureCoding.DataShards,
			"parity_shards": da.ErasureCoding.ParityShards,
			"redundancy":    fmt.Sprintf("%.0f%%", da.ErasureCoding.Redundancy*100),
		},
	}
}
