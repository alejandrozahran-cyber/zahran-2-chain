mod rpc;
mod core;
mod consensus;

use std::sync::Arc;
use tokio;
use core::{WorldState, Mempool};
use consensus::BlockProducer;

#[tokio::main]
async fn main() {
    println!("╔═══════════════════════════════════════════════════════════╗");
    println!("║                                                           ║");
    println!("║          🌌 NUSA CHAIN - FULL IMPLEMENTATION 🌌            ║");
    println!("║                                                           ║");
    println!("╚═══════════════════════════════════════════════════════════╝");
    println!();

    // Initialize state
    let state = Arc::new(WorldState::new());
    
    // Create genesis accounts
    println!("🔐 Creating genesis accounts...");
    state.create_account("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb".to_string(), 1_000_000_000_000). await;
    state.create_account("0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed".to_string(), 1_000_000_000_000). await;
    println!("✅ Genesis accounts created");

    // Initialize mempool
    let mempool = Arc::new(Mempool::new(10000));
    println!("✅ Mempool initialized (max: 10,000 txs)");

    // Start block producer (500ms blocks)
    let producer = Arc::new(BlockProducer::new(
        mempool.clone(),
        state.clone(),
        "0xValidator".to_string(),
        500, // 0.5 second block time! 
    ));
    
    println!("⚡ Block producer starting (0.5s block time)...");
    
    let producer_clone = producer.clone();
    tokio::spawn(async move {
        producer_clone.start().await;
    });

    // Start RPC server
    let server = rpc::server::RpcServer::new_with_state(state.clone(), mempool.clone());
    
    println!("🚀 NUSA Chain RPC Server starting.. .");
    println!("📡 JSON-RPC: http://0.0.0.0:8545");
    println!("🏥 Health: http://0.0.0.0:8545/health");
    println!("📊 Metrics: http://0.0.0.0:8545/metrics");
    println! ("✅ NUSA Node Ready - Processing transactions!");
    println!();

    server.run().await;
}
