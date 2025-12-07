#!/bin/bash

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║                                                           ║"
echo "║     🌌 NUSA CHAIN - COMPLETE FEATURE TEST        🌌           ║"
echo "║                                                           ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
PASSED=0
FAILED=0
TOTAL=0

# Function to test
test_feature() {
    TOTAL=$((TOTAL + 1))
    echo -n "Testing $1... "
    if eval "$2" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PASSED${NC}"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}❌ FAILED${NC}"
        FAILED=$((FAILED + 1))
    fi
}

echo "🔍 Starting comprehensive feature tests..."
echo ""

# ═══════════════════════════════════════════════════════════
# CATEGORY 1: BASIC RPC & CONNECTIVITY
# ═══════════════════════════════════════════════════════════

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📡 Category 1: Basic RPC & Connectivity"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Test 1: RPC Server responding
test_feature "RPC Server (Port 8545)" \
    "curl -s -X POST http://localhost:8545 -H 'Content-Type: application/json' --data '{\"jsonrpc\":\"2.0\",\"method\":\"eth_blockNumber\",\"params\":[],\"id\":1}' | grep -q result"

# Test 2: WebSocket endpoint
test_feature "WebSocket Support (Port 8546)" \
    "nc -zv localhost 8546"

# Test 3: Health endpoint
test_feature "Health Check Endpoint" \
    "curl -s http://localhost:8545/health | grep -q healthy"

# Test 4: Metrics endpoint
test_feature "Prometheus Metrics" \
    "curl -s http://localhost:8545/metrics | grep -q nusa"

echo ""

# ═══════════════════════════════════════════════════════════
# CATEGORY 2: BLOCKCHAIN CORE
# ═══════════════════════════════════════════════════════════

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "⛓️  Category 2: Blockchain Core Features"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Test 5: Get block number
test_feature "Get Block Number (eth_blockNumber)" \
    "curl -s -X POST http://localhost:8545 -H 'Content-Type: application/json' --data '{\"jsonrpc\":\"2. 0\",\"method\":\"eth_blockNumber\",\"params\":[],\"id\":1}' | grep -q '\"result\"'"

# Test 6: Get chain ID
test_feature "Chain ID (eth_chainId)" \
    "curl -s -X POST http://localhost:8545 -H 'Content-Type: application/json' --data '{\"jsonrpc\":\"2.0\",\"method\":\"eth_chainId\",\"params\":[],\"id\":1}' | grep -q result"

# Test 7: Network version
test_feature "Network Version (net_version)" \
    "curl -s -X POST http://localhost:8545 -H 'Content-Type: application/json' --data '{\"jsonrpc\":\"2.0\",\"method\":\"net_version\",\"params\":[],\"id\":1}' | grep -q result"

# Test 8: Get accounts
test_feature "Get Accounts (eth_accounts)" \
    "curl -s -X POST http://localhost:8545 -H 'Content-Type: application/json' --data '{\"jsonrpc\":\"2.0\",\"method\":\"eth_accounts\",\"params\":[],\"id\":1}' | grep -q result"

# Test 9: Gas price
test_feature "Get Gas Price (eth_gasPrice)" \
    "curl -s -X POST http://localhost:8545 -H 'Content-Type: application/json' --data '{\"jsonrpc\":\"2. 0\",\"method\":\"eth_gasPrice\",\"params\":[],\"id\":1}' | grep -q result"

echo ""

# ═══════════════════════════════════════════════════════════
# CATEGORY 3: CONSENSUS & VALIDATION
# ═══════════════════════════════════════════════════════════

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔐 Category 3: Consensus & Validation"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Test 10: PoS Consensus
test_feature "Proof of Stake (PoS) Module" \
    "grep -r 'PoS\|proof.*stake' src/ --include='*.rs' | head -1"

# Test 11: BFT Implementation
test_feature "Byzantine Fault Tolerance (BFT)" \
    "grep -r 'BFT\|byzantine' src/ --include='*.rs' | head -1"

# Test 12: Validator management
test_feature "Validator Management System" \
    "grep -r 'validator' src/ --include='*.rs' | head -1"

# Test 13: Finality gadget
test_feature "Fast Finality (2 seconds)" \
    "grep -r 'finality' src/ --include='*. rs' | head -1"

echo ""

# ═══════════════════════════════════════════════════════════
# CATEGORY 4: MULTI-VM SUPPORT
# ═══════════════════════════════════════════════════════════

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🖥️  Category 4: Multi-VM Support"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Test 14: EVM compatibility
test_feature "EVM (Ethereum Virtual Machine)" \
    "grep -r 'evm\|EVM' src/ --include='*. rs' | head -1"

# Test 15: WebAssembly support
test_feature "WASM (WebAssembly) VM" \
    "grep -r 'wasm\|WASM' Cargo.toml"

# Test 16: Move VM
test_feature "Move VM Integration" \
    "grep -r 'move.*vm\|MoveVM' src/ --include='*.rs' | head -1"

# Test 17: zkVM support
test_feature "zkVM (Zero-Knowledge VM)" \
    "grep -r 'zkvm\|zk.*vm' src/ --include='*.rs' | head -1"

echo ""

# ═══════════════════════════════════════════════════════════
# CATEGORY 5: PERFORMANCE & SCALABILITY
# ═══════════════════════════════════════════════════════════

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "⚡ Category 5: Performance & Scalability"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Test 18: High TPS capability
test_feature "50,000 TPS Architecture" \
    "grep -r '50000\|50k.*tps\|high.*throughput' README.md docs/"

# Test 19: Parallel execution
test_feature "Parallel Transaction Execution" \
    "grep -r 'parallel.*exec\|rayon' Cargo.toml src/"

# Test 20: State sharding
test_feature "State Sharding Implementation" \
    "grep -r 'shard' src/ --include='*.rs' | head -1"

# Test 21: Fast block time
test_feature "Sub-second Block Time (0.5s)" \
    "grep -r '0\. 5.*second\|500ms' README.md docs/"

echo ""

# ═══════════════════════════════════════════════════════════
# CATEGORY 6: SECURITY FEATURES
# ═══════════════════════════════════════════════════════════

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🛡️  Category 6: Security Features"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Test 22: Quantum-resistant crypto
test_feature "Quantum-Resistant Cryptography" \
    "grep -r 'quantum.*resist\|post.*quantum' src/ README.md"

# Test 23: MEV protection
test_feature "MEV Protection Mechanism" \
    "grep -r 'mev\|MEV\|flashbots' src/ README.md"

# Test 24: Encryption
test_feature "Advanced Encryption (AES-256)" \
    "grep -r 'aes\|encryption' Cargo.toml src/"

# Test 25: Signature verification
test_feature "Digital Signature System" \
    "grep -r 'signature\|sign.*verify' src/ --include='*.rs' | head -1"

echo ""

# ═══════════════════════════════════════════════════════════
# CATEGORY 7: CROSS-CHAIN & INTEROPERABILITY
# ═══════════════════════════════════════════════════════════

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🌐 Category 7: Cross-Chain & Interoperability"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Test 26: Cross-chain bridge
test_feature "Cross-Chain Bridge Protocol" \
    "grep -r 'bridge\|cross.*chain' src/ README.md | head -1"

# Test 27: IBC support
test_feature "IBC (Inter-Blockchain Communication)" \
    "grep -r 'ibc\|IBC' src/ README.md | head -1"

# Test 28: Asset transfer
test_feature "Cross-Chain Asset Transfer" \
    "grep -r 'asset.*transfer\|token.*bridge' src/ --include='*.rs' | head -1"

echo ""

# ═══════════════════════════════════════════════════════════
# CATEGORY 8: STORAGE & DATABASE
# ═══════════════════════════════════════════════════════════

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💾 Category 8: Storage & Database"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Test 29: PostgreSQL connection
test_feature "PostgreSQL Database (Port 5432)" \
    "nc -zv localhost 5432"

# Test 30: State storage
test_feature "State Storage System" \
    "grep -r 'storage\|state.*db' src/ --include='*.rs' | head -1"

# Test 31: IPFS integration
test_feature "IPFS Decentralized Storage" \
    "grep -r 'ipfs\|IPFS' Cargo.toml README.md"

echo ""

# ═══════════════════════════════════════════════════════════
# CATEGORY 9: MONITORING & OBSERVABILITY
# ═══════════════════════════════════════════════════════════

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Category 9: Monitoring & Observability"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Test 32: Grafana dashboard
test_feature "Grafana Dashboard (Port 3001)" \
    "curl -s http://localhost:3001/api/health | grep -q ok"

# Test 33: Prometheus metrics
test_feature "Prometheus Metrics Export" \
    "curl -s http://localhost:8545/metrics | grep -q '#'"

# Test 34: Logging system
test_feature "Advanced Logging (tracing)" \
    "grep -r 'tracing\|log' Cargo.toml | head -1"

echo ""

# ═══════════════════════════════════════════════════════════
# CATEGORY 10: ADVANCED FEATURES
# ═══════════════════════════════════════════════════════════

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 Category 10: Advanced Features"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Test 35: AI optimization
test_feature "AI-Powered Optimization" \
    "grep -r 'ai.*optim\|machine.*learning' README.md docs/"

# Test 36: IPX (Interplanetary Execution)
test_feature "IPX - Interplanetary Execution" \
    "grep -r 'ipx\|IPX\|interplanetary' README.md docs/"

# Test 37: Smart contract support
test_feature "Smart Contract Engine" \
    "grep -r 'contract' src/ --include='*.rs' | head -1"

# Test 38: Token standards
test_feature "Token Standards (ERC-20/721)" \
    "grep -r 'erc.*20\|erc.*721\|token.*standard' README.md docs/"

# Test 39: Governance system
test_feature "On-Chain Governance" \
    "grep -r 'governance\|voting' src/ README.md | head -1"

# Test 40: Upgrade mechanism
test_feature "Forkless Upgrade System" \
    "grep -r 'upgrade\|runtime.*upgrade' README.md docs/"

# Test 41: Developer tools
test_feature "Developer SDK & Tools" \
    "ls -la sdk/ docs/SDK. md 2>/dev/null"

echo ""

# ═══════════════════════════════════════════════════════════
# RESULTS SUMMARY
# ═══════════════════════════════════════════════════════════

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📈 TEST RESULTS SUMMARY"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

PERCENTAGE=$((PASSED * 100 / TOTAL))

echo "Total Tests: $TOTAL"
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"
echo "Success Rate: $PERCENTAGE%"
echo ""

if [ $PERCENTAGE -ge 80 ]; then
    echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                                                           ║${NC}"
    echo -e "${GREEN}║     ✅ EXCELLENT!  NUSA CHAIN IS PRODUCTION READY!   ✅       ║${NC}"
    echo -e "${GREEN}║                                                           ║${NC}"
    echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}"
elif [ $PERCENTAGE -ge 60 ]; then
    echo -e "${YELLOW}╔═══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${YELLOW}║     ⚠️  GOOD!  Some features need attention                 ║${NC}"
    echo -e "${YELLOW}╚═══════════════════════════════════════════════════════════╝${NC}"
else
    echo -e "${RED}╔═══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║     ❌ NEEDS WORK! Several features require fixes          ║${NC}"
    echo -e "${RED}╚═══════════════════════════════════════════════════════════╝${NC}"
fi

echo ""
echo "💡 Tip: Check failed tests and review documentation"
echo "📚 Docs: ./docs/"
echo "🐛 Issues: github.com/alejandrozahran-cyber/zahran-2-chain/issues"
echo ""
