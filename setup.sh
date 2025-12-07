#!/bin/bash

echo "🔥 Setting up NUSA Chain repository..."

# Create directory structure
mkdir -p docs
mkdir -p l1-core/cmd/nusa
mkdir -p l1-core/pkg/blockchain
mkdir -p l1-core/pkg/consensus
mkdir -p l1-core/pkg/crypto
mkdir -p l1-core/pkg/network
mkdir -p l1-core/pkg/mempool
mkdir -p l1-core/pkg/rpc
mkdir -p l2-vm/src
mkdir -p l3-ai/src
mkdir -p docker

echo "✅ Directory structure created"

# Create docker-compose.yml
cat > docker-compose.yml << 'DOCKER_COMPOSE'
version: '3.8'

services:
  nusa-l1:
    build:
      context: .
      dockerfile: docker/Dockerfile. l1
    container_name: nusa-l1-node
    ports:
      - "8080:8080"
      - "4001:4001"
    environment:
      - NODE_ENV=production
      - LOG_LEVEL=info
    networks:
      - nusa-network
    restart: unless-stopped
    volumes:
      - l1-data:/app/data

  nusa-l3:
    build:
      context: .
      dockerfile: docker/Dockerfile.l3
    container_name: nusa-l3-ai
    ports:
      - "8000:8000"
    environment:
      - PYTHONUNBUFFERED=1
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - L1_RPC_URL=http://nusa-l1:8080
    depends_on:
      - nusa-l1
      - redis
    networks:
      - nusa-network
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: nusa-redis
    ports:
      - "6379:6379"
    networks:
      - nusa-network
    restart: unless-stopped
    volumes:
      - redis-data:/data

networks:
  nusa-network:
    driver: bridge

volumes:
  l1-data:
  redis-data:
DOCKER_COMPOSE

echo "✅ docker-compose.yml created"

# Create Dockerfile. l1
cat > docker/Dockerfile.l1 << 'DOCKERFILE_L1'
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY l1-core/go.mod l1-core/go.sum ./
RUN go mod download

COPY l1-core/ . 

RUN go build -o nusa-node ./cmd/nusa

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/nusa-node .

RUN mkdir -p /app/data

EXPOSE 8080 4001

CMD ["./nusa-node"]
DOCKERFILE_L1

echo "✅ Dockerfile.l1 created"

# Create Dockerfile.l3
cat > docker/Dockerfile.l3 << 'DOCKERFILE_L3'
FROM python:3.11-slim

WORKDIR /app

RUN apt-get update && apt-get install -y gcc && rm -rf /var/lib/apt/lists/*

COPY l3-ai/requirements.txt . 

RUN pip install --no-cache-dir -r requirements.txt

COPY l3-ai/src ./src

EXPOSE 8000

CMD ["uvicorn", "src.api:app", "--host", "0.0.0. 0", "--port", "8000"]
DOCKERFILE_L3

echo "✅ Dockerfile. l3 created"

# Create l1-core/go.mod
cat > l1-core/go.mod << 'GOMOD'
module github.com/alejandrozahran-cyber/zahran-2-chain/l1-core

go 1.21

require (
	github.com/btcsuite/btcd v0.23.4
	github.com/btcsuite/btcd/btcec/v2 v2.3.2
	github.com/ethereum/go-ethereum v1.13.5
	github.com/gorilla/mux v1.8.1
	github.com/libp2p/go-libp2p v0.32.0
	github.com/multiformats/go-multiaddr v0.12.0
	github.com/tyler-smith/go-bip39 v1.1.0
)

require (
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/holiman/uint256 v1. 2.3 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)
GOMOD

echo "✅ go.mod created"

# Create go.sum (empty for now)
touch l1-core/go.sum

# Create l3-ai/requirements.txt
cat > l3-ai/requirements.txt << 'REQUIREMENTS'
fastapi==0.104.1
uvicorn[standard]==0.24.0
pydantic==2.5.0
redis==5.0.1
httpx==0.25.1
numpy==1.26.2
pandas==2.1.3
python-dotenv==1.0.0
REQUIREMENTS

echo "✅ requirements.txt created"

# Create l2-vm/Cargo.toml
cat > l2-vm/Cargo.toml << 'CARGO'
[package]
name = "nusa-vm"
version = "0.1.0"
edition = "2021"

[dependencies]
wasmtime = "14.0"
serde = { version = "1. 0", features = ["derive"] }
serde_json = "1.0"

[lib]
crate-type = ["cdylib", "rlib"]
CARGO

echo "✅ Cargo.toml created"

# Create l2-vm/src/lib.rs
cat > l2-vm/src/lib.rs << 'RUST'
//!  NUSA Chain L2 VM - WebAssembly Runtime
//!  
//! This module provides the WASM execution environment for smart contracts. 

use wasmtime::*;

pub struct NusaVM {
    engine: Engine,
}

impl NusaVM {
    pub fn new() -> Self {
        let engine = Engine::default();
        Self { engine }
    }

    pub fn execute(&self, wasm_bytes: &[u8]) -> Result<String, Box<dyn std::error::Error>> {
        // TODO: Implement WASM execution
        Ok("VM execution placeholder".to_string())
    }
}

#[no_mangle]
pub extern "C" fn execute_contract(code: *const u8, code_len: usize) -> i32 {
    // FFI entry point for Golang
    // TODO: Implement FFI bridge
    0
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_vm_creation() {
        let vm = NusaVM::new();
        assert!(true);
    }
}
RUST

echo "✅ L2 Rust placeholder created"

echo ""
echo "🎉 All files created successfully!"
echo "📁 Repository structure is ready"
echo ""
echo "Next steps:"
echo "1. Review the files"
echo "2. Run: git add ."
echo "3. Run: git commit -m 'Initial commit: NUSA Chain foundation'"
echo "4. Run: git push -u origin main"
echo ""
echo "🔥 NUSA Chain setup complete!"

