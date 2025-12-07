package validators

import (
	"fmt"
	"time"
)

// Comprehensive Validator & Network Design
// Permissionless, incentives, reputation, hardware requirements

type ValidatorNetwork struct {
	Validators           map[string]*ValidatorNode
	ActiveSet            []string
	WaitingQueue         []string
	
	// Network parameters
	MaxValidators        int
	MinStake             uint64
	PermissionlessEntry  bool
	
	// Hardware requirements
	HardwareRequirements *HardwareSpec
	
	// Incentives
	IncentiveEngine      *IncentiveEngine
	
	// Reputation
	ReputationSystem     *ReputationSystem
	
	// Governance
	GovernanceVotes      map[string]*ValidatorVote
}

type ValidatorNode struct {
	// Identity
	Address              string
	PublicKey            string
	NodeID               string
	
	// Stake
	SelfStake            uint64
	DelegatedStake       uint64
	TotalStake           uint64
	
	// Status
	Status               ValidatorStatus
	Active               bool
	Jailed               bool
	
	// Performance
	BlocksProduced       uint64
	BlocksMissed         uint64
	Uptime               float64
	AvgBlockTime         float64
	
	// Reputation
	ReputationScore      int
	
	// Commission
	CommissionRate       float64
	
	// Hardware
	Hardware             *ValidatorHardware
	
	// Network
	NetworkInfo          *NetworkInfo
	
	// Rewards
	TotalRewards         uint64
	
	// Slashing
	SlashingHistory      []SlashEvent
	
	// Metadata
	Name                 string
	Website              string
	Description          string
	JoinedAt             time.Time
	LastActiveAt         time.Time
}

type ValidatorStatus string

const (
	StatusActive      ValidatorStatus = "active"
	StatusWaiting     ValidatorStatus = "waiting"
	StatusJailed      ValidatorStatus = "jailed"
	StatusUnbonding   ValidatorStatus = "unbonding"
	StatusSlashed     ValidatorStatus = "slashed"
)

type HardwareSpec struct {
	MinCPU               int     // Cores
	MinRAM               int     // GB
	MinDisk              int     // GB
	MinBandwidth         int     // Mbps
	RecommendedCPU       int
	RecommendedRAM       int
	RecommendedDisk      int
	RecommendedBandwidth int
}

type ValidatorHardware struct {
	CPU                  int
	RAM                  int
	Disk                 int
	Bandwidth            int
	MeetsRequirements    bool
	Score                int  // Hardware quality score
}

type NetworkInfo struct {
	IPAddress            string
	Port                 int
	P2PPort              int
	RPCPort              int
	Latency              time.Duration
	PeerCount            int
}

type IncentiveEngine struct {
	// Block rewards
	BlockRewardPool      uint64
	
	// Performance bonuses
	UptimeBonus          float64  // +20% for >99.9% uptime
	FastBlockBonus       float64  // +10% for fast blocks
	
	// Penalties
	MissedBlockPenalty   uint64
	DowntimePenalty      float64
}

type ReputationSystem struct {
	Scores               map[string]*ReputationScore
	DecayRate            float64  // Reputation decays over time
}

type ReputationScore struct {
	Validator            string
	Score                int
	Components           map[string]int
	LastUpdated          time.Time
}

type SlashEvent struct {
	Reason               string
	Amount               uint64
	Timestamp            time.Time
	BlockHeight          uint64
}

type ValidatorVote struct {
	Validator            string
	ProposalID           string
	Vote                 bool
	Weight               uint64
	Timestamp            time. Time
}

func NewValidatorNetwork() *ValidatorNetwork {
	return &ValidatorNetwork{
		Validators:          make(map[string]*ValidatorNode),
		ActiveSet:           make([]string, 0),
		WaitingQueue:        make([]string, 0),
		
		MaxValidators:       100,  // Top 100 by stake
		MinStake:            100_000 * 1e8,  // 100K NUSA
		PermissionlessEntry: true,
		
		HardwareRequirements: &HardwareSpec{
			MinCPU:               4,
			MinRAM:               16,
			MinDisk:              500,
			MinBandwidth:         100,
			RecommendedCPU:       8,
			RecommendedRAM:       32,
			RecommendedDisk:      1000,
			RecommendedBandwidth: 1000,
		},
		
		IncentiveEngine:  NewIncentiveEngine(),
		ReputationSystem: NewReputationSystem(),
		GovernanceVotes:  make(map[string]*ValidatorVote),
	}
}

func NewIncentiveEngine() *IncentiveEngine {
	return &IncentiveEngine{
		BlockRewardPool:    1_000_000 * 1e8,
		UptimeBonus:        0.20,  // 20%
		FastBlockBonus:     0.10,  // 10%
		MissedBlockPenalty: 100 * 1e8,
		DowntimePenalty:    0.05,  // 5%
	}
}

func NewReputationSystem() *ReputationSystem {
	return &ReputationSystem{
		Scores:    make(map[string]*ReputationScore),
		DecayRate: 0.01,  // 1% per week
	}
}

// Register validator (permissionless!)
func (vn *ValidatorNetwork) RegisterValidator(
	address, publicKey, name string,
	stake uint64,
	hardware *ValidatorHardware,
	network *NetworkInfo,
	commissionRate float64,
) error {
	// 1. Check minimum stake
	if stake < vn.MinStake {
		return fmt.Errorf("insufficient stake: %d < %d", stake, vn. MinStake)
	}
	
	// 2. Check hardware requirements
	if ! vn.meetsHardwareRequirements(hardware) {
		return fmt. Errorf("hardware does not meet minimum requirements")
	}
	
	// 3.  Validate commission rate
	if commissionRate < 0 || commissionRate > 0.20 {  // Max 20%
		return fmt.Errorf("commission rate must be 0-20%%")
	}
	
	// 4. Create validator node
	validator := &ValidatorNode{
		Address:         address,
		PublicKey:       publicKey,
		NodeID:          generateNodeID(address),
		
		SelfStake:       stake,
		DelegatedStake:  0,
		TotalStake:      stake,
		
		Status:          StatusWaiting,
		Active:          false,
		Jailed:          false,
		
		BlocksProduced:  0,
		BlocksMissed:    0,
		Uptime:          100.0,
		AvgBlockTime:    2.0,
		
		ReputationScore: 100,
		
		CommissionRate:  commissionRate,
		
		Hardware:        hardware,
		NetworkInfo:     network,
		
		TotalRewards:    0,
		
		SlashingHistory: make([]SlashEvent, 0),
		
		Name:            name,
		JoinedAt:        time.Now(),
		LastActiveAt:    time.Now(),
	}
	
	vn.Validators[address] = validator
	
	// 5. Add to waiting queue or active set
	vn.updateValidatorSet()
	
	fmt.Printf("✅ Validator registered: %s (stake: %d NUSA)\n", name, stake/1e8)
	
	// 6. Initialize reputation
	vn.ReputationSystem.initializeReputation(address)
	
	return nil
}

// Check hardware requirements
func (vn *ValidatorNetwork) meetsHardwareRequirements(hw *ValidatorHardware) bool {
	req := vn.HardwareRequirements
	
	meets := hw.CPU >= req.MinCPU &&
		hw.RAM >= req.MinRAM &&
		hw.Disk >= req.MinDisk &&
		hw. Bandwidth >= req.MinBandwidth
	
	hw.MeetsRequirements = meets
	
	// Calculate hardware score
	hw.Score = vn.calculateHardwareScore(hw)
	
	return meets
}

func (vn *ValidatorNetwork) calculateHardwareScore(hw *ValidatorHardware) int {
	req := vn.HardwareRequirements
	
	score := 0
	
	// CPU score
	score += (hw.CPU * 100) / req.RecommendedCPU
	
	// RAM score
	score += (hw. RAM * 100) / req. RecommendedRAM
	
	// Disk score
	score += (hw.Disk * 100) / req.RecommendedDisk
	
	// Bandwidth score
	score += (hw.Bandwidth * 100) / req.RecommendedBandwidth
	
	return score / 4  // Average
}

// Update active validator set
func (vn *ValidatorNetwork) updateValidatorSet() {
	// Sort validators by total stake
	var sortedValidators []string
	
	for addr, val := range vn.Validators {
		if val.Status != StatusJailed && val.Status != StatusSlashed {
			sortedValidators = append(sortedValidators, addr)
		}
	}
	
	// Sort by stake (simplified - use heap in production)
	// Top N become active validators
	
	activeCount := len(sortedValidators)
	if activeCount > vn.MaxValidators {
		activeCount = vn.MaxValidators
	}
	
	vn.ActiveSet = sortedValidators[:activeCount]
	vn.WaitingQueue = sortedValidators[activeCount:]
	
	// Update status
	for _, addr := range vn.ActiveSet {
		vn. Validators[addr].Status = StatusActive
		vn. Validators[addr].Active = true
	}
	
	for _, addr := range vn.WaitingQueue {
		vn.Validators[addr].Status = StatusWaiting
		vn.Validators[addr].Active = false
	}
	
	fmt.Printf("🔄 Validator set updated: %d active, %d waiting\n",
		len(vn.ActiveSet), len(vn.WaitingQueue))
}

// Record block production
func (vn *ValidatorNetwork) RecordBlockProduction(validator string, blockTime float64) {
	val, exists := vn.Validators[validator]
	if !exists {
		return
	}
	
	val.BlocksProduced++
	val.LastActiveAt = time.Now()
	
	// Update average block time
	val.AvgBlockTime = (val.AvgBlockTime + blockTime) / 2. 0
	
	// Calculate reward
	reward := vn.IncentiveEngine.calculateBlockReward(val)
	val.TotalRewards += reward
	
	// Update reputation
	vn.ReputationSystem.updateReputation(validator, "block_produced", 1)
	
	fmt.Printf("⛏️ Block produced by %s (time: %.2fs, reward: %d NUSA)\n",
		val.Name, blockTime, reward/1e8)
}

// Record missed block
func (vn *ValidatorNetwork) RecordMissedBlock(validator string) {
	val, exists := vn.Validators[validator]
	if !exists {
		return
	}
	
	val.BlocksMissed++
	
	// Apply penalty
	penalty := vn.IncentiveEngine.MissedBlockPenalty
	if val.TotalRewards >= penalty {
		val.TotalRewards -= penalty
	}
	
	// Update reputation (negative)
	vn.ReputationSystem.updateReputation(validator, "block_missed", -5)
	
	// Update uptime
	totalBlocks := val.BlocksProduced + val. BlocksMissed
	val. Uptime = (float64(val.BlocksProduced) / float64(totalBlocks)) * 100
	
	fmt.Printf("⚠️ Block missed by %s (uptime: %.2f%%)\n", val.Name, val.Uptime)
	
	// Jail if too many misses
	if val.BlocksMissed > 100 && val.Uptime < 95.0 {
		vn.JailValidator(validator, "excessive_missed_blocks")
	}
}

// Jail validator
func (vn *ValidatorNetwork) JailValidator(validator, reason string) {
	val, exists := vn.Validators[validator]
	if !exists {
		return
	}
	
	val.Status = StatusJailed
	val. Active = false
	val.Jailed = true
	
	// Remove from active set
	vn. updateValidatorSet()
	
	fmt.Printf("🔒 Validator jailed: %s (reason: %s)\n", val.Name, reason)
}

// Unjail validator (after fixing issues)
func (vn *ValidatorNetwork) UnjailValidator(validator string) {
	val, exists := vn.Validators[validator]
	if !exists || ! val.Jailed {
		return
	}
	
	val. Jailed = false
	val. Status = StatusWaiting
	
	vn.updateValidatorSet()
	
	fmt.Printf("🔓 Validator unjailed: %s\n", val. Name)
}

// Delegate stake to validator
func (vn *ValidatorNetwork) Delegate(delegator, validator string, amount uint64) error {
	val, exists := vn.Validators[validator]
	if !exists {
		return fmt.Errorf("validator not found")
	}
	
	val.DelegatedStake += amount
	val.TotalStake += amount
	
	// Update validator set (may change ranking)
	vn.updateValidatorSet()
	
	fmt.Printf("🤝 Delegated: %d NUSA from %s to %s (total: %d)\n",
		amount/1e8, delegator, val.Name, val.TotalStake/1e8)
	
	return nil
}

// Calculate block reward with bonuses
func (ie *IncentiveEngine) calculateBlockReward(val *ValidatorNode) uint64 {
	baseReward := uint64(10 * 1e8)  // 10 NUSA base
	
	totalReward := float64(baseReward)
	
	// Uptime bonus
	if val.Uptime >= 99.9 {
		totalReward += totalReward * ie.UptimeBonus
	}
	
	// Fast block bonus
	if val. AvgBlockTime < 2.0 {
		totalReward += totalReward * ie.FastBlockBonus
	}
	
	return uint64(totalReward)
}

// Reputation system
func (rs *ReputationSystem) initializeReputation(validator string) {
	rs.Scores[validator] = &ReputationScore{
		Validator: validator,
		Score:     100,
		Components: map[string]int{
			"uptime":          100,
			"block_production": 100,
			"governance":      100,
			"community":       100,
		},
		LastUpdated: time.Now(),
	}
}

func (rs *ReputationSystem) updateReputation(validator, component string, change int) {
	score, exists := rs.Scores[validator]
	if !exists {
		return
	}
	
	// Update component
	if current, ok := score.Components[component]; ok {
		score.Components[component] = current + change
		
		// Bounds check
		if score.Components[component] < 0 {
			score.Components[component] = 0
		} else if score.Components[component] > 100 {
			score.Components[component] = 100
		}
	}
	
	// Recalculate total score
	total := 0
	for _, val := range score.Components {
		total += val
	}
	score.Score = total / len(score.Components)
	
	score.LastUpdated = time.Now()
}

// Get validator network stats
func (vn *ValidatorNetwork) GetStats() map[string]interface{} {
	totalStake := uint64(0)
	avgUptime := 0.0
	
	for _, val := range vn.ActiveSet {
		v := vn.Validators[val]
		totalStake += v.TotalStake
		avgUptime += v.Uptime
	}
	
	if len(vn.ActiveSet) > 0 {
		avgUptime /= float64(len(vn.ActiveSet))
	}
	
	return map[string]interface{}{
		"total_validators":  len(vn. Validators),
		"active_validators": len(vn.ActiveSet),
		"waiting_validators": len(vn.WaitingQueue),
		"total_stake":       totalStake / 1e8,
		"avg_uptime":        fmt.Sprintf("%.2f%%", avgUptime),
		"min_stake":         vn.MinStake / 1e8,
		"permissionless":    vn.PermissionlessEntry,
		"hardware_requirements": vn.HardwareRequirements,
	}
}

func generateNodeID(address string) string {
	return fmt.Sprintf("node_%s", address[:16])
}
