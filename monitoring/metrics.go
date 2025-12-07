package monitoring

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Metrics & Monitoring System
// Prometheus-compatible metrics for validators, nodes, consensus

type MetricsCollector struct {
	mu sync.RWMutex
	
	// Blockchain metrics
	BlockHeight          uint64
	TotalTransactions    uint64
	PendingTransactions  uint64
	AverageBlockTime     float64
	CurrentTPS           float64
	PeakTPS              float64
	
	// Consensus metrics
	ValidatorCount       int
	ActiveValidators     int
	MissedBlocks         uint64
	ConsensusRounds      uint64
	AverageRoundTime     float64
	
	// Network metrics
	PeerCount            int
	InboundConnections   int
	OutboundConnections  int
	NetworkBandwidth     uint64 // bytes/sec
	MessagesSent         uint64
	MessagesReceived     uint64
	
	// Node health
	CPUUsage             float64
	MemoryUsage          float64
	DiskUsage            float64
	Uptime               time.Duration
	SyncStatus           string
	
	// Gas & fees
	AverageGasPrice      uint64
	TotalFeesCollected   uint64
	TotalFeesBurned      uint64
	
	// State metrics
	TotalAccounts        uint64
	TotalContracts       uint64
	StateSize            uint64 // bytes
	
	// Performance
	BlockProcessingTime  float64 // milliseconds
	TxValidationTime     float64
	StateUpdateTime      float64
	
	StartTime            time.Time
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		StartTime:  time.Now(),
		SyncStatus: "syncing",
	}
}

// Update block metrics
func (mc *MetricsCollector) RecordBlock(blockNumber uint64, txCount int, processingTime float64) {
	mc.mu.Lock()
	defer mc. mu.Unlock()
	
	mc.BlockHeight = blockNumber
	mc.TotalTransactions += uint64(txCount)
	mc.BlockProcessingTime = processingTime
	
	// Calculate TPS
	mc.CurrentTPS = float64(txCount) / (processingTime / 1000. 0)
	if mc.CurrentTPS > mc. PeakTPS {
		mc.PeakTPS = mc.CurrentTPS
	}
	
	fmt.Printf("📊 Block #%d | TXs: %d | TPS: %.2f | Time: %.2fms\n",
		blockNumber, txCount, mc.CurrentTPS, processingTime)
}

// Update consensus metrics
func (mc *MetricsCollector) RecordConsensusRound(validators int, roundTime float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	
	mc.ConsensusRounds++
	mc.ActiveValidators = validators
	mc. AverageRoundTime = (mc.AverageRoundTime + roundTime) / 2.0
}

// Record validator missed block
func (mc *MetricsCollector) RecordMissedBlock(validator string) {
	mc.mu.Lock()
	defer mc. mu.Unlock()
	
	mc.MissedBlocks++
	fmt.Printf("⚠️ Validator %s missed block\n", validator)
}

// Update network metrics
func (mc *MetricsCollector) UpdateNetworkMetrics(peers, inbound, outbound int) {
	mc.mu.Lock()
	defer mc.mu. Unlock()
	
	mc. PeerCount = peers
	mc.InboundConnections = inbound
	mc.OutboundConnections = outbound
}

// Update node health
func (mc *MetricsCollector) UpdateNodeHealth(cpu, memory, disk float64) {
	mc.mu. Lock()
	defer mc.mu.Unlock()
	
	mc.CPUUsage = cpu
	mc.MemoryUsage = memory
	mc.DiskUsage = disk
	mc.Uptime = time.Since(mc.StartTime)
}

// Update gas metrics
func (mc *MetricsCollector) RecordGasMetrics(avgGasPrice, feesCollected, feesBurned uint64) {
	mc.mu.Lock()
	defer mc.mu. Unlock()
	
	mc. AverageGasPrice = avgGasPrice
	mc.TotalFeesCollected += feesCollected
	mc. TotalFeesBurned += feesBurned
}

// Prometheus metrics endpoint
func (mc *MetricsCollector) ServeMetrics() {
	http.HandleFunc("/metrics", mc.handleMetrics)
	
	fmt.Println("📊 Metrics server running on :9090")
	http.ListenAndServe(":9090", nil)
}

func (mc *MetricsCollector) handleMetrics(w http. ResponseWriter, r *http.Request) {
	mc.mu. RLock()
	defer mc.mu.RUnlock()
	
	// Prometheus format
	metrics := fmt.Sprintf(`# HELP nusa_block_height Current block height
# TYPE nusa_block_height gauge
nusa_block_height %d

# HELP nusa_total_transactions Total transactions processed
# TYPE nusa_total_transactions counter
nusa_total_transactions %d

# HELP nusa_pending_transactions Pending transactions in mempool
# TYPE nusa_pending_transactions gauge
nusa_pending_transactions %d

# HELP nusa_current_tps Current transactions per second
# TYPE nusa_current_tps gauge
nusa_current_tps %. 2f

# HELP nusa_peak_tps Peak TPS achieved
# TYPE nusa_peak_tps gauge
nusa_peak_tps %.2f

# HELP nusa_validator_count Total validators
# TYPE nusa_validator_count gauge
nusa_validator_count %d

# HELP nusa_active_validators Active validators
# TYPE nusa_active_validators gauge
nusa_active_validators %d

# HELP nusa_missed_blocks Total missed blocks
# TYPE nusa_missed_blocks counter
nusa_missed_blocks %d

# HELP nusa_consensus_rounds Total consensus rounds
# TYPE nusa_consensus_rounds counter
nusa_consensus_rounds %d

# HELP nusa_peer_count Connected peers
# TYPE nusa_peer_count gauge
nusa_peer_count %d

# HELP nusa_cpu_usage CPU usage percentage
# TYPE nusa_cpu_usage gauge
nusa_cpu_usage %.2f

# HELP nusa_memory_usage Memory usage percentage
# TYPE nusa_memory_usage gauge
nusa_memory_usage %. 2f

# HELP nusa_disk_usage Disk usage percentage
# TYPE nusa_disk_usage gauge
nusa_disk_usage %.2f

# HELP nusa_uptime_seconds Node uptime in seconds
# TYPE nusa_uptime_seconds counter
nusa_uptime_seconds %. 0f

# HELP nusa_average_gas_price Average gas price
# TYPE nusa_average_gas_price gauge
nusa_average_gas_price %d

# HELP nusa_fees_collected Total fees collected
# TYPE nusa_fees_collected counter
nusa_fees_collected %d

# HELP nusa_fees_burned Total fees burned
# TYPE nusa_fees_burned counter
nusa_fees_burned %d

# HELP nusa_total_accounts Total accounts
# TYPE nusa_total_accounts gauge
nusa_total_accounts %d

# HELP nusa_state_size_bytes State size in bytes
# TYPE nusa_state_size_bytes gauge
nusa_state_size_bytes %d

# HELP nusa_block_processing_time_ms Block processing time in milliseconds
# TYPE nusa_block_processing_time_ms gauge
nusa_block_processing_time_ms %. 2f
`,
		mc.BlockHeight,
		mc.TotalTransactions,
		mc.PendingTransactions,
		mc.CurrentTPS,
		mc.PeakTPS,
		mc.ValidatorCount,
		mc.ActiveValidators,
		mc.MissedBlocks,
		mc.ConsensusRounds,
		mc.PeerCount,
		mc.CPUUsage,
		mc.MemoryUsage,
		mc. DiskUsage,
		mc.Uptime. Seconds(),
		mc.AverageGasPrice,
		mc.TotalFeesCollected,
		mc.TotalFeesBurned,
		mc.TotalAccounts,
		mc.StateSize,
		mc.BlockProcessingTime,
	)
	
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(metrics))
}

// Get dashboard summary
func (mc *MetricsCollector) GetDashboard() map[string]interface{} {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	return map[string]interface{}{
		"blockchain": map[string]interface{}{
			"block_height":         mc.BlockHeight,
			"total_transactions":   mc.TotalTransactions,
			"pending_transactions": mc.PendingTransactions,
			"current_tps":          fmt.Sprintf("%.2f", mc.CurrentTPS),
			"peak_tps":             fmt.Sprintf("%.2f", mc.PeakTPS),
			"avg_block_time":       fmt.Sprintf("%.2fs", mc.AverageBlockTime),
		},
		"consensus": map[string]interface{}{
			"validators":       mc.ValidatorCount,
			"active":           mc.ActiveValidators,
			"missed_blocks":    mc.MissedBlocks,
			"consensus_rounds": mc. ConsensusRounds,
			"avg_round_time":   fmt. Sprintf("%.2fms", mc.AverageRoundTime),
		},
		"network": map[string]interface{}{
			"peers":    mc.PeerCount,
			"inbound":  mc.InboundConnections,
			"outbound": mc.OutboundConnections,
		},
		"node_health": map[string]interface{}{
			"cpu":    fmt.Sprintf("%.1f%%", mc.CPUUsage),
			"memory": fmt.Sprintf("%.1f%%", mc.MemoryUsage),
			"disk":   fmt.Sprintf("%.1f%%", mc.DiskUsage),
			"uptime": mc. Uptime.String(),
			"status": mc.SyncStatus,
		},
		"economics": map[string]interface{}{
			"avg_gas_price":   mc.AverageGasPrice,
			"fees_collected":  mc.TotalFeesCollected,
			"fees_burned":     mc.TotalFeesBurned,
		},
	}
}

// Alert system
type AlertManager struct {
	alerts []Alert
	mu     sync. Mutex
}

type Alert struct {
	Level     string    // critical, warning, info
	Message   string
	Timestamp time.Time
	Resolved  bool
}

func NewAlertManager() *AlertManager {
	return &AlertManager{
		alerts: make([]Alert, 0),
	}
}

func (am *AlertManager) Trigger(level, message string) {
	am. mu.Lock()
	defer am.mu.Unlock()
	
	alert := Alert{
		Level:     level,
		Message:   message,
		Timestamp: time. Now(),
		Resolved:  false,
	}
	
	am.alerts = append(am.alerts, alert)
	
	emoji := "ℹ️"
	if level == "warning" {
		emoji = "⚠️"
	} else if level == "critical" {
		emoji = "🚨"
	}
	
	fmt.Printf("%s ALERT [%s]: %s\n", emoji, level, message)
}

func (am *AlertManager) GetActiveAlerts() []Alert {
	am.mu.Lock()
	defer am.mu.Unlock()
	
	active := make([]Alert, 0)
	for _, alert := range am.alerts {
		if ! alert.Resolved {
			active = append(active, alert)
		}
	}
	
	return active
}
