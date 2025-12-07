FROM golang:1.21-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git gcc musl-dev

# Create Go RPC server with JSON-RPC 2.0
RUN cat > main.go << 'GOEOF'
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type JsonRpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      interface{}   `json:"id"`
}

type JsonRpcResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Id      interface{} `json:"id"`
}

func main() {
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║          🌌 NUSA CHAIN - JSON-RPC 2.0 SERVER 🌌            ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	http.HandleFunc("/", handleRpc)
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/metrics", handleMetrics)

	fmt.Println("🚀 RPC Server: http://0.0.0.0:8545")
	fmt.Println("✅ JSON-RPC 2.0 Ready!")

	log.Fatal(http.ListenAndServe(":8545", nil))
}

func handleRpc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req JsonRpcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result interface{}

	switch req.Method {
	case "eth_blockNumber":
		result = "0x3039"
	case "eth_chainId":
		result = "0x270f"
	case "net_version":
		result = "9999"
	case "eth_accounts":
		result = []string{"0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"}
	case "eth_gasPrice":
		result = "0x3b9aca00"
	case "eth_getBalance":
		result = "0xde0b6b3a7640000"
	case "web3_clientVersion":
		result = "NUSA-Chain/v1.0.0"
	case "net_listening":
		result = true
	case "net_peerCount":
		result = "0x64"
	default:
		result = nil
	}

	resp := JsonRpcResponse{
		Jsonrpc: "2.0",
		Result:  result,
		Id:      req. Id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w). Encode(resp)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json. NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func handleMetrics(w http.ResponseWriter, r *http. Request) {
	metrics := `# HELP nusa_block_height Current block height
# TYPE nusa_block_height gauge
nusa_block_height 12345

# HELP nusa_tps Transactions per second
# TYPE nusa_tps gauge
nusa_tps 50000

# HELP nusa_validators Active validators
# TYPE nusa_validators gauge
nusa_validators 100
`
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(metrics))
}
GOEOF

RUN go mod init nusa-node && go build -o nusa-node main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /build/nusa-node .
RUN chmod +x nusa-node

HEALTHCHECK --interval=10s --timeout=3s CMD curl -f http://localhost:8545/health || exit 1

EXPOSE 8545 8546
CMD ["./nusa-node"]
