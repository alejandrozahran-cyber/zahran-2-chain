#!/bin/bash
PASSED=0
FAILED=0

check() {
    printf "Testing %-50s" "$1..."
    if eval "$2" > /dev/null 2>&1; then
        echo " ✅ PASS"
        ((PASSED++))
    else
        echo " ❌ FAIL"
        ((FAILED++))
    fi
}

echo "╔═══════════════════════════════════════════════════════════╗"
echo "║     🔍 NUSA CHAIN - ACCURATE TEST 🔍                       ║"
echo "╚═══════════════════════════════════════════════════════════╝"

check "1. RPC" "curl -sf http://localhost:8545/health | grep -q healthy"
check "2. WebSocket" "nc -zv localhost 8546 2>&1 | grep -q succeeded"
check "3. Health" "curl -sf http://localhost:8545/health | grep -q healthy"
check "4.  Metrics" "curl -sf http://localhost:8545/metrics | grep -q nusa"

check "5. eth_blockNumber" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2. 0\",\"method\":\"eth_blockNumber\",\"params\":[],\"id\":1}' | grep -q 0x"
check "6. eth_chainId" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2. 0\",\"method\":\"eth_chainId\",\"params\":[],\"id\":1}' | grep -q 0x4e555341"
check "7. net_version" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"net_version\",\"params\":[],\"id\":1}' | grep -q 1313376900"
check "8. eth_accounts" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"eth_accounts\",\"params\":[],\"id\":1}' | grep -q 0x742d35"
check "9.  eth_gasPrice" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2. 0\",\"method\":\"eth_gasPrice\",\"params\":[],\"id\":1}' | grep -q 0x"

check "10. PoS" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2. 0\",\"method\":\"nusa_posInfo\",\"params\":[],\"id\":1}' | grep -q PoS"
check "11.  BFT" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_bftInfo\",\"params\":[],\"id\":1}' | grep -q BFT"
check "12.  Validators" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_validators\",\"params\":[],\"id\":1}' | grep -q total"
check "13.  Finality" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2. 0\",\"method\":\"nusa_finality\",\"params\":[],\"id\":1}' | grep -q 2s"

check "14. EVM" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_evmInfo\",\"params\":[],\"id\":1}' | grep -q operational"
check "15. WASM" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2. 0\",\"method\":\"nusa_wasmInfo\",\"params\":[],\"id\":1}' | grep -q operational"
check "16. Move" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_moveInfo\",\"params\":[],\"id\":1}' | grep -q operational"
check "17. zkVM" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_zkInfo\",\"params\":[],\"id\":1}' | grep -q operational"

check "18. TPS" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_tpsInfo\",\"params\":[],\"id\":1}' | grep -q 50000"
check "19.  Parallel" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_parallelExecution\",\"params\":[],\"id\":1}' | grep -q enabled"
check "20. Sharding" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2. 0\",\"method\":\"nusa_shardingInfo\",\"params\":[],\"id\":1}' | grep -q shards"
check "21. BlockTime" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_blockTime\",\"params\":[],\"id\":1}' | grep -q 0.5s"

check "22. Quantum" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_quantumInfo\",\"params\":[],\"id\":1}' | grep -q Dilithium"
check "23.  MEV" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2. 0\",\"method\":\"nusa_mevProtection\",\"params\":[],\"id\":1}' | grep -q enabled"
check "24. AES" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2. 0\",\"method\":\"nusa_encryptionInfo\",\"params\":[],\"id\":1}' | grep -q AES"
check "25. Signature" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_signatureInfo\",\"params\":[],\"id\":1}' | grep -q ECDSA"

check "26. Bridge" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_bridgeInfo\",\"params\":[],\"id\":1}' | grep -q ETH"
check "27. IBC" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2. 0\",\"method\":\"nusa_ibcInfo\",\"params\":[],\"id\":1}' | grep -q IBC"
check "28. CrossChain" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_crossChainTransfer\",\"params\":[],\"id\":1}' | grep -q supported"

check "29. PostgreSQL" "docker ps | grep -q postgres"
check "30.  Storage" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_storageInfo\",\"params\":[],\"id\":1}' | grep -q PostgreSQL"
check "31.  IPFS" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_ipfsInfo\",\"params\":[],\"id\":1}' | grep -q enabled"

check "32. Grafana" "docker ps | grep -q grafana"
check "33. Metrics" "curl -sf http://localhost:8545/metrics | grep -q nusa"
check "34. Logging" "docker logs nusa-node 2>&1 | grep -q ✅"

check "35. AI" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_aiInfo\",\"params\":[],\"id\":1}' | grep -q ML"
check "36. IPX" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_ipxInfo\",\"params\":[],\"id\":1}' | grep -q interplanetary"
check "37. Contracts" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_contractInfo\",\"params\":[],\"id\":1}' | grep -q operational"
check "38. Tokens" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_tokenInfo\",\"params\":[],\"id\":1}' | grep -q erc20"
check "39.  Governance" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_governanceInfo\",\"params\":[],\"id\":1}' | grep -q voting"
check "40. Upgrade" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_upgradeInfo\",\"params\":[],\"id\":1}' | grep -q forkless"
check "41. SDK" "curl -sf -XPOST http://localhost:8545 -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"method\":\"nusa_sdkInfo\",\"params\":[],\"id\":1}' | grep -q Rust"

TOTAL=$((PASSED + FAILED))
PERCENT=$((PASSED * 100 / TOTAL))
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Results: $PASSED/$TOTAL passed ($PERCENT%)"
if [ $PERCENT -ge 70 ]; then echo "🎉 EXCELLENT! "; elif [ $PERCENT -ge 50 ]; then echo "✅ GOOD! "; else echo "⚠️  NEEDS WORK"; fi
