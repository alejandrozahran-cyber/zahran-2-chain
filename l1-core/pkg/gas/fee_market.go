package gas

import (
	"fmt"
	"math"
)

// EIP-1559 Style Fee Market
// Dynamic base fee + priority fee (tip)

type FeeMarket struct {
	BaseFee            uint64  // Current base fee (adjusted per block)
	MaxBaseFee         uint64  // Max base fee cap
	MinBaseFee         uint64  // Min base fee floor
	TargetGasPerBlock  uint64  // Target: 50% full blocks
	MaxGasPerBlock     uint64  // Maximum gas per block
	BaseFeeChangeDenom uint64  // Denominator for base fee change (8 = 12. 5% max change)
	BurnedFees         uint64  // Total fees burned (deflationary)
}

type Transaction struct {
	From             string
	To               string
	Value            uint64
	GasLimit         uint64 // Max gas willing to use
	MaxFeePerGas     uint64 // Max total fee per gas
	MaxPriorityFee   uint64 // Max tip to validator
	Data             []byte
	GasUsed          uint64
	EffectiveFee     uint64
}

func NewFeeMarket() *FeeMarket {
	return &FeeMarket{
		BaseFee:            1000,      // 0.00001 NUSA
		MaxBaseFee:         1000000,   // 0.01 NUSA
		MinBaseFee:         100,       // 0.000001 NUSA
		TargetGasPerBlock:  15000000,  // 15M gas (50% of max)
		MaxGasPerBlock:     30000000,  // 30M gas max
		BaseFeeChangeDenom: 8,         // 12.5% max change per block
		BurnedFees:         0,
	}
}

// Calculate effective fee for transaction
func (fm *FeeMarket) CalculateEffectiveFee(tx *Transaction) uint64 {
	// Effective fee = min(maxFeePerGas, baseFee + maxPriorityFee)
	maxFee := tx.MaxFeePerGas
	targetFee := fm.BaseFee + tx.MaxPriorityFee

	effectiveFee := uint64(math.Min(float64(maxFee), float64(targetFee)))

	return effectiveFee
}

// Validate transaction fee
func (fm *FeeMarket) ValidateTransaction(tx *Transaction) (bool, string) {
	// Check if maxFeePerGas >= baseFee
	if tx.MaxFeePerGas < fm.BaseFee {
		return false, fmt.Sprintf("maxFeePerGas (%d) < baseFee (%d)", tx.MaxFeePerGas, fm.BaseFee)
	}

	// Check gas limit
	if tx.GasLimit > fm.MaxGasPerBlock {
		return false, fmt.Sprintf("gasLimit (%d) exceeds max (%d)", tx.GasLimit, fm.MaxGasPerBlock)
	}

	// Check if user has enough balance for max fee
	maxCost := tx.GasLimit * tx.MaxFeePerGas
	// TODO: Check actual balance

	return true, "valid"
}

// Execute transaction and calculate fees
func (fm *FeeMarket) ExecuteTransaction(tx *Transaction) (uint64, uint64, error) {
	// Validate first
	valid, reason := fm.ValidateTransaction(tx)
	if !valid {
		return 0, 0, fmt.Errorf("invalid transaction: %s", reason)
	}

	// Calculate effective fee
	effectiveFee := fm. CalculateEffectiveFee(tx)

	// Simulate gas usage (production: actual execution)
	tx.GasUsed = tx.GasLimit * 7 / 10 // Assume 70% usage

	// Calculate fees
	totalFee := tx.GasUsed * effectiveFee
	baseFeeAmount := tx.GasUsed * fm.BaseFee
	priorityFeeAmount := totalFee - baseFeeAmount

	// Burn base fee (deflationary!)
	fm.BurnedFees += baseFeeAmount

	// Priority fee goes to validator
	validatorReward := priorityFeeAmount

	fmt.Printf("💸 TX executed | Gas: %d | Base: %d | Priority: %d | Burned: %d | Validator: %d\n",
		tx.GasUsed, fm.BaseFee, tx.MaxPriorityFee, baseFeeAmount, validatorReward)

	return baseFeeAmount, validatorReward, nil
}

// Adjust base fee after block (EIP-1559 algorithm)
func (fm *FeeMarket) AdjustBaseFee(blockGasUsed uint64) {
	if blockGasUsed == fm.TargetGasPerBlock {
		// Perfect utilization - no change
		return
	}

	oldBaseFee := fm.BaseFee

	if blockGasUsed > fm.TargetGasPerBlock {
		// Block more than 50% full - increase base fee
		gasUsedDelta := blockGasUsed - fm.TargetGasPerBlock
		baseFeeChange := fm.BaseFee * gasUsedDelta / fm.TargetGasPerBlock / fm.BaseFeeChangeDenom

		if baseFeeChange < 1 {
			baseFeeChange = 1
		}

		fm.BaseFee += baseFeeChange

		// Cap at max
		if fm.BaseFee > fm. MaxBaseFee {
			fm.BaseFee = fm.MaxBaseFee
		}

	} else {
		// Block less than 50% full - decrease base fee
		gasUsedDelta := fm. TargetGasPerBlock - blockGasUsed
		baseFeeChange := fm.BaseFee * gasUsedDelta / fm.TargetGasPerBlock / fm.BaseFeeChangeDenom

		if baseFeeChange < 1 {
			baseFeeChange = 1
		}

		if fm.BaseFee > baseFeeChange {
			fm. BaseFee -= baseFeeChange
		}

		// Floor at min
		if fm.BaseFee < fm.MinBaseFee {
			fm.BaseFee = fm.MinBaseFee
		}
	}

	changePercent := float64(fm.BaseFee-oldBaseFee) / float64(oldBaseFee) * 100

	fmt.Printf("⛽ Base fee adjusted: %d → %d (%.2f%%) | Block utilization: %.1f%%\n",
		oldBaseFee, fm.BaseFee, changePercent, float64(blockGasUsed)/float64(fm.TargetGasPerBlock)*50)
}

// Estimate gas for operation
func (fm *FeeMarket) EstimateGas(operationType string) uint64 {
	gasEstimates := map[string]uint64{
		"transfer":        21000,
		"token_transfer":  65000,
		"swap":            150000,
		"nft_mint":        100000,
		"contract_deploy": 500000,
		"dao_vote":        80000,
	}

	if gas, exists := gasEstimates[operationType]; exists {
		return gas
	}

	return 21000 // Default
}

// Get recommended fees for users
func (fm *FeeMarket) GetRecommendedFees() map[string]uint64 {
	return map[string]uint64{
		"slow":     fm.BaseFee + 1000,       // Just above base fee
		"standard": fm.BaseFee + 5000,       // Normal priority
		"fast":     fm. BaseFee + 20000,      // High priority
		"instant":  fm.BaseFee + 100000,     // Immediate inclusion
	}
}

// Get fee market stats
func (fm *FeeMarket) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"base_fee":         fm.BaseFee,
		"burned_fees":      fm.BurnedFees,
		"target_gas":       fm.TargetGasPerBlock,
		"max_gas":          fm.MaxGasPerBlock,
		"recommended_fees": fm.GetRecommendedFees(),
	}
}
