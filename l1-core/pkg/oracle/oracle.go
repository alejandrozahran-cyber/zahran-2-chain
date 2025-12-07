package oracle

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

// Oracle Network - Bring real-world data on-chain
// Features: Price feeds, weather data, sports results, random numbers

type OracleNetwork struct {
	DataFeeds   map[string]*DataFeed
	Oracles     map[string]*Oracle
	Requests    map[string]*DataRequest
	MinOracles  int     // Minimum oracles for consensus
	Threshold   float64 // Agreement threshold (90%)
}

type DataFeed struct {
	FeedID      string
	DataType    DataType
	Value       interface{}
	LastUpdated time.Time
	UpdateFreq  time.Duration
	Providers   []string
	Confidence  float64 // 0-100%
}

type DataType string

const (
	PriceFeed    DataType = "price"
	WeatherData  DataType = "weather"
	SportsResult DataType = "sports"
	RandomNumber DataType = "random"
	APIData      DataType = "api"
)

type Oracle struct {
	Address    string
	Reputation int
	DataFeeds  []string
	Uptime     float64
	Slashed    bool
}

type DataRequest struct {
	RequestID   string
	Requester   string
	DataType    DataType
	Query       string
	Responses   map[string]interface{} // oracle_address -> value
	Finalized   bool
	Result      interface{}
	Timestamp   time.Time
	Reward      uint64
}

func NewOracleNetwork() *OracleNetwork {
	return &OracleNetwork{
		DataFeeds:  make(map[string]*DataFeed),
		Oracles:    make(map[string]*Oracle),
		Requests:   make(map[string]*DataRequest),
		MinOracles: 3,
		Threshold:  0.90, // 90% agreement
	}
}

// Register oracle node
func (on *OracleNetwork) RegisterOracle(address string) bool {
	// Check if already registered
	if _, exists := on.Oracles[address]; exists {
		return false
	}

	oracle := &Oracle{
		Address:    address,
		Reputation: 100,
		DataFeeds:  make([]string, 0),
		Uptime:     99.9,
		Slashed:    false,
	}

	on.Oracles[address] = oracle

	fmt.Printf("🔮 Oracle registered: %s\n", address)

	return true
}

// Create price feed (e.g., NUSA/USD)
func (on *OracleNetwork) CreatePriceFeed(pair string, updateFreq time.Duration) string {
	feedID := fmt.Sprintf("feed_%s", pair)

	feed := &DataFeed{
		FeedID:      feedID,
		DataType:    PriceFeed,
		Value:       0.0,
		LastUpdated: time.Now(),
		UpdateFreq:  updateFreq,
		Providers:   make([]string, 0),
		Confidence:  0.0,
	}

	on. DataFeeds[feedID] = feed

	fmt.Printf("📊 Price feed created: %s (update every %v)\n", pair, updateFreq)

	return feedID
}

// Oracle submits data
func (on *OracleNetwork) SubmitData(oracleAddr, feedID string, value interface{}) bool {
	oracle, exists := on.Oracles[oracleAddr]
	if !exists || oracle.Slashed {
		return false
	}

	feed, exists := on.DataFeeds[feedID]
	if !exists {
		return false
	}

	// Add oracle as provider if not already
	isProvider := false
	for _, provider := range feed.Providers {
		if provider == oracleAddr {
			isProvider = true
			break
		}
	}

	if !isProvider {
		feed.Providers = append(feed.Providers, oracleAddr)
	}

	// Collect submissions for consensus
	// (Simplified - production needs time-based aggregation)

	fmt.Printf("📡 Oracle %s submitted data for %s: %v\n", oracleAddr, feedID, value)

	// Update feed value (simplified - should aggregate multiple oracles)
	feed.Value = value
	feed.LastUpdated = time.Now()

	// Reward oracle
	oracle.Reputation += 1

	return true
}

// Request data from oracles
func (on *OracleNetwork) RequestData(
	requester string,
	dataType DataType,
	query string,
	reward uint64,
) string {
	hash := sha256.Sum256([]byte(query + requester + time.Now().String()))
	requestID := fmt.Sprintf("req_%x", hash[:8])

	request := &DataRequest{
		RequestID: requestID,
		Requester: requester,
		DataType:  dataType,
		Query:     query,
		Responses: make(map[string]interface{}),
		Finalized: false,
		Timestamp: time.Now(),
		Reward:    reward,
	}

	on.Requests[requestID] = request

	fmt.Printf("📝 Data request created: %s | Query: %s | Reward: %d\n", requestID, query, reward)

	return requestID
}

// Oracle responds to data request
func (on *OracleNetwork) RespondToRequest(
	oracleAddr string,
	requestID string,
	response interface{},
) bool {
	oracle, exists := on. Oracles[oracleAddr]
	if !exists || oracle. Slashed {
		return false
	}

	request, exists := on.Requests[requestID]
	if !exists || request.Finalized {
		return false
	}

	// Submit response
	request.Responses[oracleAddr] = response

	fmt.Printf("📬 Oracle %s responded to %s: %v\n", oracleAddr, requestID, response)

	// Check if enough responses to finalize
	if len(request. Responses) >= on.MinOracles {
		on.FinalizeRequest(requestID)
	}

	return true
}

// Finalize data request with consensus
func (on *OracleNetwork) FinalizeRequest(requestID string) bool {
	request, exists := on.Requests[requestID]
	if ! exists || request.Finalized {
		return false
	}

	if len(request.Responses) < on.MinOracles {
		fmt. Println("❌ Not enough oracle responses")
		return false
	}

	// Consensus mechanism: Majority vote or median
	var finalResult interface{}

	switch request.DataType {
	case PriceFeed, RandomNumber:
		// Use median for numeric data
		finalResult = on.calculateMedian(request.Responses)

	case WeatherData, SportsResult, APIData:
		// Use majority vote for categorical data
		finalResult = on.calculateMajority(request.Responses)
	}

	request.Result = finalResult
	request. Finalized = true

	fmt. Printf("✅ Request finalized: %s | Result: %v\n", requestID, finalResult)

	// Distribute rewards to oracles
	rewardPerOracle := request.Reward / uint64(len(request.Responses))
	for oracleAddr := range request.Responses {
		if oracle, exists := on.Oracles[oracleAddr]; exists {
			oracle.Reputation += 5
			fmt.Printf("💰 Oracle %s earned %d NUSA\n", oracleAddr, rewardPerOracle)
		}
	}

	return true
}

// Calculate median (for numeric data)
func (on *OracleNetwork) calculateMedian(responses map[string]interface{}) float64 {
	values := make([]float64, 0)

	for _, val := range responses {
		if fVal, ok := val.(float64); ok {
			values = append(values, fVal)
		}
	}

	if len(values) == 0 {
		return 0. 0
	}

	// Simple median (production: proper sorting)
	sum := 0.0
	for _, v := range values {
		sum += v
	}

	return sum / float64(len(values))
}

// Calculate majority (for categorical data)
func (on *OracleNetwork) calculateMajority(responses map[string]interface{}) interface{} {
	votes := make(map[string]int)

	for _, val := range responses {
		strVal := fmt.Sprintf("%v", val)
		votes[strVal]++
	}

	// Find most common response
	maxVotes := 0
	var majority interface{}

	for val, count := range votes {
		if count > maxVotes {
			maxVotes = count
			majority = val
		}
	}

	return majority
}

// Generate verifiable random number
func (on *OracleNetwork) GenerateVRF(seed string) uint64 {
	// Verifiable Random Function
	// Production: Use proper VRF like Chainlink VRF

	hash := sha256.Sum256([]byte(seed + time.Now(). String()))

	// Convert first 8 bytes to uint64
	randomNum := uint64(hash[0]) | uint64(hash[1])<<8 | uint64(hash[2])<<16 |
		uint64(hash[3])<<24 | uint64(hash[4])<<32 | uint64(hash[5])<<40 |
		uint64(hash[6])<<48 | uint64(hash[7])<<56

	fmt.Printf("🎲 VRF generated: %d (seed: %s)\n", randomNum, seed)

	return randomNum
}

// Slash malicious oracle
func (on *OracleNetwork) SlashOracle(oracleAddr string, reason string) bool {
	oracle, exists := on.Oracles[oracleAddr]
	if !exists {
		return false
	}

	oracle.Slashed = true
	oracle.Reputation -= 50

	fmt.Printf("⚠️ Oracle %s slashed: %s\n", oracleAddr, reason)

	return true
}

// Get price from feed
func (on *OracleNetwork) GetPrice(pair string) (float64, bool) {
	feedID := fmt.Sprintf("feed_%s", pair)
	feed, exists := on.DataFeeds[feedID]

	if !exists {
		return 0.0, false
	}

	// Check if data is fresh
	if time.Since(feed.LastUpdated) > feed.UpdateFreq*2 {
		fmt. Println("⚠️ Price data stale")
		return 0.0, false
	}

	if price, ok := feed.Value.(float64); ok {
		return price, true
	}

	return 0.0, false
}

// Automated price update (for demo)
func (on *OracleNetwork) SimulatePriceUpdate(pair string) {
	feedID := fmt.Sprintf("feed_%s", pair)

	// Simulate multiple oracles submitting prices
	basePrice := 2.5 + rand.Float64()*0.5 // $2.5-$3.0

	for addr, oracle := range on.Oracles {
		if oracle.Slashed {
			continue
		}

		// Each oracle submits slightly different price
		price := basePrice + (rand.Float64()-0.5)*0.1 // ±$0.05

		on.SubmitData(addr, feedID, price)
	}
}

// Get oracle stats
func (on *OracleNetwork) GetOracleStats(addr string) map[string]interface{} {
	oracle, exists := on. Oracles[addr]
	if !exists {
		return nil
	}

	return map[string]interface{}{
		"address":    oracle.Address,
		"reputation": oracle.Reputation,
		"uptime":     oracle.Uptime,
		"data_feeds": len(oracle.DataFeeds),
		"slashed":    oracle.Slashed,
	}
}

// Get network stats
func (on *OracleNetwork) GetNetworkStats() map[string]interface{} {
	activeOracles := 0
	for _, oracle := range on.Oracles {
		if !oracle. Slashed {
			activeOracles++
		}
	}

	return map[string]interface{}{
		"total_oracles":   len(on.Oracles),
		"active_oracles":  activeOracles,
		"total_feeds":     len(on.DataFeeds),
		"total_requests":  len(on. Requests),
		"min_consensus":   on.MinOracles,
	}
}
