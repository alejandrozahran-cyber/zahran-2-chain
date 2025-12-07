package dao

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// DAO Governance System
type DAO struct {
	Proposals       map[string]*Proposal
	Members         map[string]*Member
	TreasuryBalance uint64
	QuorumRequired  float64 // 10% of members must vote
}

type Proposal struct {
	ID          string
	Title       string
	Description string
	Proposer    string
	Category    ProposalCategory
	VotesFor    uint64
	VotesAgainst uint64
	Status      ProposalStatus
	CreatedAt   time.Time
	Deadline    time.Time
	Executed    bool
}

type ProposalCategory string

const (
	TreasurySpend   ProposalCategory = "treasury_spend"
	ParameterChange ProposalCategory = "parameter_change"
	UpgradeProtocol ProposalCategory = "upgrade_protocol"
	AddValidator    ProposalCategory = "add_validator"
	EmergencyAction ProposalCategory = "emergency"
)

type ProposalStatus string

const (
	Active   ProposalStatus = "active"
	Passed   ProposalStatus = "passed"
	Rejected ProposalStatus = "rejected"
	Executed ProposalStatus = "executed"
)

type Member struct {
	Address       string
	VotingPower   uint64 // Based on NUSA holdings + reputation
	Reputation    int
	ProposalsVoted int
	DelegatedTo   string // Can delegate voting power
}

func NewDAO() *DAO {
	return &DAO{
		Proposals:       make(map[string]*Proposal),
		Members:         make(map[string]*Member),
		TreasuryBalance: 0,
		QuorumRequired:  0.10, // 10%
	}
}

// Create new proposal
func (d *DAO) CreateProposal(
	title, description, proposer string,
	category ProposalCategory,
) string {
	// Check if proposer is member
	member, exists := d.Members[proposer]
	if !exists {
		return ""
	}

	// Require minimum voting power to propose
	if member.VotingPower < 1000 {
		fmt.Println("❌ Insufficient voting power to propose")
		return ""
	}

	// Generate proposal ID
	hash := sha256.Sum256([]byte(title + proposer + time.Now().String()))
	proposalID := fmt.Sprintf("prop_%x", hash[:8])

	proposal := &Proposal{
		ID:           proposalID,
		Title:        title,
		Description:  description,
		Proposer:     proposer,
		Category:     category,
		VotesFor:     0,
		VotesAgainst: 0,
		Status:       Active,
		CreatedAt:    time.Now(),
		Deadline:     time.Now().Add(7 * 24 * time. Hour), // 7 days voting period
		Executed:     false,
	}

	d.Proposals[proposalID] = proposal

	fmt.Printf("📝 Proposal created: %s by %s\n", title, proposer)

	return proposalID
}

// Vote on proposal
func (d *DAO) Vote(proposalID, voter string, support bool) bool {
	proposal, exists := d.Proposals[proposalID]
	if !exists {
		return false
	}

	member, exists := d.Members[voter]
	if !exists {
		return false
	}

	// Check if voting period still active
	if time.Now().After(proposal.Deadline) {
		fmt.Println("❌ Voting period ended")
		return false
	}

	// Apply voting power
	votingPower := member.VotingPower

	// Check if delegated
	if member.DelegatedTo != "" {
		fmt.Println("❌ Voting power delegated to another member")
		return false
	}

	if support {
		proposal.VotesFor += votingPower
		fmt.Printf("✅ %s voted FOR with power %d\n", voter, votingPower)
	} else {
		proposal.VotesAgainst += votingPower
		fmt.Printf("❌ %s voted AGAINST with power %d\n", voter, votingPower)
	}

	member.ProposalsVoted++
	member. Reputation += 1 // Reward participation

	return true
}

// Finalize proposal after voting period
func (d *DAO) FinalizeProposal(proposalID string) bool {
	proposal, exists := d. Proposals[proposalID]
	if !exists {
		return false
	}

	// Check if deadline passed
	if time.Now().Before(proposal.Deadline) {
		fmt.Println("⏳ Voting still in progress")
		return false
	}

	// Calculate total votes
	totalVotes := proposal. VotesFor + proposal.VotesAgainst
	totalVotingPower := d.getTotalVotingPower()

	// Check quorum
	quorum := float64(totalVotes) / float64(totalVotingPower)
	if quorum < d.QuorumRequired {
		proposal.Status = Rejected
		fmt.Printf("❌ Proposal failed: Quorum not reached (%.2f%% < %.2f%%)\n", 
			quorum*100, d.QuorumRequired*100)
		return false
	}

	// Check if passed (simple majority)
	if proposal.VotesFor > proposal.VotesAgainst {
		proposal.Status = Passed
		fmt.Printf("✅ Proposal PASSED!  (%d for, %d against)\n", 
			proposal.VotesFor, proposal.VotesAgainst)
		
		// Auto-execute if applicable
		d.ExecuteProposal(proposalID)
		
		return true
	} else {
		proposal.Status = Rejected
		fmt.Printf("❌ Proposal REJECTED (%d for, %d against)\n", 
			proposal.VotesFor, proposal.VotesAgainst)
		return false
	}
}

// Execute approved proposal
func (d *DAO) ExecuteProposal(proposalID string) bool {
	proposal, exists := d.Proposals[proposalID]
	if !exists || proposal.Status != Passed {
		return false
	}

	if proposal.Executed {
		fmt.Println("⚠️ Proposal already executed")
		return false
	}

	// Execute based on category
	switch proposal.Category {
	case TreasurySpend:
		fmt.Printf("💰 Executing treasury spend: %s\n", proposal.Title)
		// TODO: Transfer funds
		
	case ParameterChange:
		fmt.Printf("⚙️ Changing protocol parameters: %s\n", proposal.Title)
		// TODO: Update parameters
		
	case UpgradeProtocol:
		fmt.Printf("🔄 Upgrading protocol: %s\n", proposal.Title)
		// TODO: Deploy upgrade
		
	case AddValidator:
		fmt.Printf("✅ Adding new validator: %s\n", proposal.Title)
		// TODO: Add validator
	}

	proposal.Executed = true
	proposal.Status = Executed

	return true
}

// Delegate voting power
func (d *DAO) Delegate(from, to string) bool {
	fromMember, exists := d.Members[from]
	if !exists {
		return false
	}

	_, exists = d.Members[to]
	if !exists {
		return false
	}

	fromMember.DelegatedTo = to
	fmt.Printf("🤝 %s delegated voting power to %s\n", from, to)

	return true
}

// Add member to DAO
func (d *DAO) AddMember(address string, votingPower uint64) {
	d.Members[address] = &Member{
		Address:       address,
		VotingPower:   votingPower,
		Reputation:    0,
		ProposalsVoted: 0,
		DelegatedTo:   "",
	}
}

func (d *DAO) getTotalVotingPower() uint64 {
	total := uint64(0)
	for _, member := range d.Members {
		total += member.VotingPower
	}
	return total
}

// Get DAO stats
func (d *DAO) GetStats() map[string]interface{} {
	activeProposals := 0
	for _, p := range d.Proposals {
		if p.Status == Active {
			activeProposals++
		}
	}

	return map[string]interface{}{
		"total_members":      len(d.Members),
		"active_proposals":   activeProposals,
		"treasury_balance":   d.TreasuryBalance,
		"total_voting_power": d.getTotalVotingPower(),
	}
}
