use std::sync::Arc;
use tokio::sync::RwLock;
use tokio::time::{interval, Duration};
use crate::core::{Block, Mempool, TransactionExecutor, WorldState};

pub struct BlockProducer {
    mempool: Arc<Mempool>,
    executor: Arc<TransactionExecutor>,
    current_block: Arc<RwLock<u64>>,
    last_hash: Arc<RwLock<String>>,
    validator_address: String,
    block_time_ms: u64,
}

impl BlockProducer {
    pub fn new(
        mempool: Arc<Mempool>,
        state: Arc<WorldState>,
        validator_address: String,
        block_time_ms: u64,
    ) -> Self {
        BlockProducer {
            mempool,
            executor: Arc::new(TransactionExecutor::new(state)),
            current_block: Arc::new(RwLock::new(0)),
            last_hash: Arc::new(RwLock::new(String::from("0x0"))),
            validator_address,
            block_time_ms,
        }
    }

    pub async fn start(&self) {
        let mut interval = interval(Duration::from_millis(self.block_time_ms));
        
        loop {
            interval.tick().await;
            self.produce_block().await;
        }
    }

    async fn produce_block(&self) {
        let block_number = {
            let mut current = self.current_block. write().await;
            *current += 1;
            *current
        };

        let parent_hash = self.last_hash.read().await.clone();

        // Get transactions from mempool
        let transactions = self.mempool.get_transactions(1000).await;
        
        if transactions.is_empty() {
            // No transactions, skip block
            return;
        }

        // Create new block
        let mut block = Block::new(block_number, parent_hash, self.validator_address.clone());

        // Execute transactions
        let receipts = self.executor.execute_batch(transactions. clone(), block_number).await;

        // Add successful transactions to block
        for (tx, receipt) in transactions.iter().zip(receipts.iter()) {
            if receipt.status {
                block. add_transaction(tx.clone());
            }
        }

        // Update last hash
        {
            let mut last_hash = self.last_hash.write().await;
            *last_hash = block.hash. clone();
        }

        println!(
            "⛓️  Block #{} produced | {} txs | Hash: {}",
            block. number,
            block.transactions. len(),
            &block.hash[..16]
        );
    }

    pub async fn get_current_block_number(&self) -> u64 {
        *self.current_block.read(). await
    }
}
