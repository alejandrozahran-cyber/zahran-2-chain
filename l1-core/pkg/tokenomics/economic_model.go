package tokenomics

import (
	"fmt"
	"math"
	"time"
)

// Comprehensive Tokenomics & Economic Model
// Token supply, inflation, rewards, slashing, velocity control

type TokenomicsEngine struct {
	TotalSupply          uint64
	CirculatingSupply    uint64
	MaxSupply            uint64
	InitialPrice         float64
	CurrentPrice         float64
	
	// Inflation control
	InflationRate        float64
	TargetInflationRate  float64
	MaxInflationRate     float64
	MinInflationRate     float64
	
	// Supply schedule
	SupplySchedule       *SupplySchedule
	
	// Staking
	StakingPool          *StakingPool
	
	// Burning
	TotalBurned          uint64
	BurnRate             float64
	
	// Velocity control
	VelocityTracker      *VelocityTracker
	
	// UBI
	UBIPool              uint64
	DailyUBIPerPerson    uint64
	
	// Treasury
	Treasury             uint64
}

type SupplySchedule struct {
	Phases               []SupplyPhase
	CurrentPhase         int
	NextHalvingBlock     uint64
}

type SupplyPhase struct {
	Name                 string
	StartBlock           uint64
	EndBlock             uint64
	BlockReward          uint64
	InflationRate        float64
	Description          string
}

type StakingPool struct {
	TotalStaked          uint64
	StakingRatio         float64  // % of circulating supply staked
	MinStake             uint64
	UnbondingPeriod      time.Duration
	
	// Rewards
	BaseAPY              float64
	CurrentAPY           float64
	MaxAPY               float64
	RewardsDistributed   uint64
	
	// Validators
	Validators           map[string]*Validator
}

type Validator struct {
	Address              string
	Stake                uint64
	Delegated            uint64
	Commission           float64
	Reputation           int
	BlocksProduced       uint64
	BlocksMissed         uint64
	Slashed              bool
	SlashedAmount        uint64
	JoinedAt             time.Time
}

type VelocityTracker struct {
	DailyTransactions    uint64
	DailyVolume          uint64
	TokenVelocity        float64  // transactions per token per day
	TargetVelocity       float64
	VelocityMultiplier   float64  // Affects inflation
}

func NewTokenomicsEngine() *TokenomicsEngine {
	return &TokenomicsEngine{
		TotalSupply:         1_000_000_000 * 1e8,  // 1 billion NUSA
		CirculatingSupply:   300_000_000 * 1e8,    // 30% initial circulation
		MaxSupply:           10_000_000_000 * 1e8, // 10 billion cap
		InitialPrice:        2.5,
		CurrentPrice:        2.5,
		
		InflationRate:       0.05,  // 5% per year
		TargetInflationRate: 0.03,  // 3% target
		MaxInflationRate:    0.10,  // 10% max
		MinInflationRate:    0.01,  // 1% min
		
		SupplySchedule:      NewSupplySchedule(),
		StakingPool:         NewStakingPool(),
		
		TotalBurned:         0,
		BurnRate:            0.01,  // 1% of fees burned
		
		VelocityTracker:     NewVelocityTracker(),
		
		UBIPool:             50_000_000 * 1e8,  // 5% for UBI
		DailyUBIPerPerson:   10 * 1e8,
		
		Treasury:            100_000_000 * 1e8, // 10% treasury
	}
}

func NewSupplySchedule() *SupplySchedule {
	return &SupplySchedule{
		Phases: []SupplyPhase{
			{
				Name:          "Genesis Phase",
				StartBlock:    0,
				EndBlock:      10_512_000,  // ~2 years
				BlockReward:   10 * 1e8,     // 10 NUSA per block
				InflationRate: 0.05,
				Description:   "Initial high rewards for network bootstrapping",
			},
			{
				Name:          "Growth Phase",
				StartBlock:    10_512_000,
				EndBlock:      21_024_000,  // Years 2-4
				BlockReward:   5 * 1e8,      // 5 NUSA per block (50% reduction)
				InflationRate: 0.03,
				Description:   "Moderate rewards for sustained growth",
			},
			{
				Name:          "Maturity Phase",
				StartBlock:    21_024_000,
				EndBlock:      52_560_000,  // Years 4-10
				BlockReward:   2. 5 * 1e8,    // 2.5 NUSA per block
				InflationRate: 0.02,
				Description:   "Low inflation, mature network",
			},
			{
				Name:          "Steady State",
				StartBlock:    52_560_000,
				EndBlock:      math.MaxUint64,
				BlockReward:   1 * 1e8,      // 1 NUSA per block
				InflationRate: 0.01,
				Description:   "Minimal inflation, fee-driven security",
			},
		},
		CurrentPhase:     0,
		NextHalvingBlock: 10_512_000,
	}
}

func NewStakingPool() *StakingPool {
	return &StakingPool{
		TotalStaked:        0,
		StakingRatio:       0,
		MinStake:           1000 * 1e8,  // 1000 NUSA minimum
		UnbondingPeriod:    21 * 24 * time.Hour,  // 21 days
		
		BaseAPY:            0.08,   // 8% base APY
		CurrentAPY:         0.08,
		MaxAPY:             0. 15,   // 15% max APY
		RewardsDistributed: 0,
		
		Validators:         make(map[string]*Validator),
	}
}

func NewVelocityTracker() *VelocityTracker {
	return &VelocityTracker{
		DailyTransactions:  0,
		DailyVolume:        0,
		TokenVelocity:      1.0,
		TargetVelocity:     2.0,  // Healthy velocity = 2x per day
		VelocityMultiplier: 1.0,
	}
}

// Calculate block reward based on current phase
func (te *TokenomicsEngine) GetBlockReward(blockNumber uint64) uint64 {
	phase := te.SupplySchedule.GetCurrentPhase(blockNumber)
	return phase.BlockReward
}

// Dynamic inflation adjustment based on economic conditions
func (te *TokenomicsEngine) AdjustInflation(blockNumber uint64) {
	// Factors affecting inflation:
	// 1. Staking ratio (higher staking = lower inflation)
	// 2.  Token velocity (higher velocity = lower inflation needed)
	// 3. Network growth (more users = can sustain higher inflation)
	
	stakingRatio := te.StakingPool.StakingRatio
	velocity := te.VelocityTracker.TokenVelocity
	
	// Base inflation from supply schedule
	phase := te.SupplySchedule.GetCurrentPhase(blockNumber)
	baseInflation := phase.InflationRate
	
	// Adjust based on staking (more staking = reduce inflation)
	stakingAdjustment := 1.0 - (stakingRatio * 0.2)  // Max 20% reduction
	
	// Adjust based on velocity (higher velocity = reduce inflation)
	velocityAdjustment := 1.0
	if velocity > te.VelocityTracker. TargetVelocity {
		velocityAdjustment = 0.9  // Reduce inflation by 10%
	} else if velocity < te.VelocityTracker.TargetVelocity * 0.5 {
		velocityAdjustment = 1.1  // Increase inflation by 10%
	}
	
	// Calculate new inflation
	newInflation := baseInflation * stakingAdjustment * velocityAdjustment
	
	// Apply bounds
	if newInflation > te.MaxInflationRate {
		newInflation = te.MaxInflationRate
	} else if newInflation < te.MinInflationRate {
		newInflation = te. MinInflationRate
	}
	
	oldInflation := te.InflationRate
	te.InflationRate = newInflation
	
	fmt.Printf("📊 Inflation adjusted: %. 2f%% → %.2f%% (staking: %.1f%%, velocity: %.2f)\n",
		oldInflation*100, newInflation*100, stakingRatio*100, velocity)
}

// Stake tokens
func (te *TokenomicsEngine) Stake(validator string, amount uint64) error {
	if amount < te.StakingPool.MinStake {
		return fmt.Errorf("minimum stake: %d", te.StakingPool.MinStake)
	}
	
	// Get or create validator
	val, exists := te.StakingPool.Validators[validator]
	if !exists {
		val = &Validator{
			Address:        validator,
			Stake:          0,
			Delegated:      0,
			Commission:     0. 10,  // 10% default commission
			Reputation:     100,
			BlocksProduced: 0,
			BlocksMissed:   0,
			Slashed:        false,
			SlashedAmount:  0,
			JoinedAt:       time.Now(),
		}
		te.StakingPool. Validators[validator] = val
	}
	
	// Update staking
	val.Stake += amount
	te.StakingPool.TotalStaked += amount
	te.CirculatingSupply -= amount  // Remove from circulation
	
	// Update staking ratio
	te.StakingPool.StakingRatio = float64(te.StakingPool.TotalStaked) / float64(te.CirculatingSupply + te.StakingPool.TotalStaked)
	
	// Adjust APY based on staking ratio
	te.adjustStakingAPY()
	
	fmt.Printf("🔒 Staked: %d NUSA by %s (total staked: %d, ratio: %.1f%%)\n",
		amount/1e8, validator, te. StakingPool.TotalStaked/1e8, te.StakingPool.StakingRatio*100)
	
	return nil
}

// Adjust staking APY dynamically
func (te *TokenomicsEngine) adjustStakingAPY() {
	// Higher staking ratio = lower APY (supply/demand)
	// Target: 40-60% staking ratio
	
	ratio := te.StakingPool. StakingRatio
	
	if ratio < 0.40 {
		// Too little staked, increase APY to attract stakers
		te.StakingPool.CurrentAPY = te.StakingPool.MaxAPY
	} else if ratio > 0.60 {
		// Too much staked, reduce APY
		te.StakingPool. CurrentAPY = te.StakingPool.BaseAPY * 0.7
	} else {
		// Optimal range
		te.StakingPool.CurrentAPY = te.StakingPool.BaseAPY
	}
}

// Slash validator for misbehavior
func (te *TokenomicsEngine) SlashValidator(validator string, reason string, percentage float64) uint64 {
	val, exists := te.StakingPool.Validators[validator]
	if !exists {
		return 0
	}
	
	slashAmount := uint64(float64(val.Stake) * percentage)
	
	val. Stake -= slashAmount
	val. Slashed = true
	val.SlashedAmount += slashAmount
	val. Reputation -= 50
	
	te.StakingPool.TotalStaked -= slashAmount
	
	// Burn 50% of slashed amount, rest to treasury
	burnAmount := slashAmount / 2
	treasuryAmount := slashAmount - burnAmount
	
	te. Burn(burnAmount, "validator_slash")
	te.Treasury += treasuryAmount
	
	fmt.Printf("⚠️ VALIDATOR SLASHED: %s (%. 0f%%) for %s | Burned: %d | Treasury: %d\n",
		validator, percentage*100, reason, burnAmount/1e8, treasuryAmount/1e8)
	
	return slashAmount
}

// Burn tokens (deflationary mechanism)
func (te *TokenomicsEngine) Burn(amount uint64, reason string) {
	te.TotalBurned += amount
	te.TotalSupply -= amount
	
	fmt.Printf("🔥 Burned: %d NUSA (reason: %s) | Total burned: %d\n",
		amount/1e8, reason, te. TotalBurned/1e8)
}

// Calculate token velocity
func (te *TokenomicsEngine) UpdateVelocity(dailyTxs, dailyVolume uint64) {
	te.VelocityTracker.DailyTransactions = dailyTxs
	te.VelocityTracker.DailyVolume = dailyVolume
	
	// Velocity = daily volume / circulating supply
	te.VelocityTracker.TokenVelocity = float64(dailyVolume) / float64(te.CirculatingSupply)
	
	fmt.Printf("💨 Token Velocity: %.2f (target: %.2f)\n",
		te.VelocityTracker.TokenVelocity, te.VelocityTracker.TargetVelocity)
}

// Get tokenomics stats
func (te *TokenomicsEngine) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"supply": map[string]interface{}{
			"total":       te.TotalSupply / 1e8,
			"circulating": te.CirculatingSupply / 1e8,
			"max":         te.MaxSupply / 1e8,
			"burned":      te.TotalBurned / 1e8,
		},
		"inflation": map[string]interface{}{
			"current":     fmt.Sprintf("%.2f%%", te.InflationRate*100),
			"target":      fmt.Sprintf("%.2f%%", te.TargetInflationRate*100),
		},
		"staking": map[string]interface{}{
			"total_staked":  te.StakingPool. TotalStaked / 1e8,
			"staking_ratio": fmt.Sprintf("%.1f%%", te.StakingPool. StakingRatio*100),
			"current_apy":   fmt.Sprintf("%.1f%%", te.StakingPool.CurrentAPY*100),
			"validators":    len(te.StakingPool.Validators),
		},
		"velocity": map[string]interface{}{
			"current":          fmt.Sprintf("%.2f", te.VelocityTracker.TokenVelocity),
			"target":           fmt. Sprintf("%.2f", te. VelocityTracker. TargetVelocity),
		},
		"ubi": map[string]interface{}{
			"pool":        te.UBIPool / 1e8,
			"daily_per_person": te.DailyUBIPerPerson / 1e8,
		},
		"treasury":    te.Treasury / 1e8,
		"price":       fmt.Sprintf("$%.2f", te.CurrentPrice),
		"market_cap":  fmt.Sprintf("$%.2fB", float64(te.CirculatingSupply/1e8)*te.CurrentPrice/1e9),
	}
}

// Supply schedule helpers
func (ss *SupplySchedule) GetCurrentPhase(blockNumber uint64) SupplyPhase {
	for i, phase := range ss.Phases {
		if blockNumber >= phase.StartBlock && blockNumber < phase.EndBlock {
			ss.CurrentPhase = i
			return phase
		}
	}
	return ss.Phases[len(ss.Phases)-1]  // Return last phase if beyond all
}
