#!/bin/bash

echo "🧪 Testing sub-second block time..."

BLOCK1=$(curl -s -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  | grep -o '"result":"[^"]*"' | cut -d'"' -f4)

sleep 0.6

BLOCK2=$(curl -s -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  | grep -o '"result":"[^"]*"' | cut -d'"' -f4)

NUM1=$((16#${BLOCK1#0x}))
NUM2=$((16#${BLOCK2#0x}))

if [ $NUM2 -gt $NUM1 ]; then
    echo "✅ PASSED - Block increased from $NUM1 to $NUM2 in 0.6s"
    exit 0
else
    echo "❌ FAILED - Block didn't increase"
    exit 1
fi
