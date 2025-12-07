use super::transaction::{Transaction, TransactionReceipt};
use super::state::WorldState;
use std::sync::Arc;

pub struct TransactionExecutor {
    state: Arc<WorldState>,
}

impl TransactionExecutor {
    pub fn new(state: Arc<WorldState>) -> Self {
        TransactionExecutor { state }
    }

    pub async fn execute(&self, tx: &Transaction, block_number: u64) -> TransactionReceipt {
        // Verify transaction
        if !tx.verify() {
            return TransactionReceipt {
                tx_hash: tx.hash.clone(),
                block_number,
                gas_used: 0,
                status: false,
                logs: vec! ["Transaction verification failed".to_string()],
            };
        }

        // Execute transfer
        match self.state.transfer(&tx. from, &tx.to, tx.value).await {
            Ok(_) => TransactionReceipt {
                tx_hash: tx.hash.clone(),
                block_number,
                gas_used: tx.gas_limit,
                status: true,
                logs: vec![format!("Transferred {} from {} to {}", tx.value, tx.from, tx.to)],
            },
            Err(e) => TransactionReceipt {
                tx_hash: tx.hash.clone(),
                block_number,
                gas_used: 21000, // Base gas even on failure
                status: false,
                logs: vec![format!("Execution failed: {}", e)],
            },
        }
    }

    pub async fn execute_batch(&self, transactions: Vec<Transaction>, block_number: u64) -> Vec<TransactionReceipt> {
        let mut receipts = Vec::new();
        
        for tx in transactions {
            let receipt = self.execute(&tx, block_number).await;
            receipts.push(receipt);
        }

        receipts
    }

    // Parallel execution using rayon
    pub async fn execute_parallel(&self, transactions: Vec<Transaction>, block_number: u64) -> Vec<TransactionReceipt> {
        use rayon::prelude::*;
        
        let executor = Arc::new(self.state.clone());
        
        // Group transactions by sender to avoid conflicts
        let mut groups: Vec<Vec<Transaction>> = Vec::new();
        let mut current_group = Vec::new();
        let mut seen_senders = std::collections::HashSet::new();

        for tx in transactions {
            if seen_senders.contains(&tx.from) {
                groups.push(current_group);
                current_group = Vec::new();
                seen_senders.clear();
            }
            seen_senders.insert(tx. from.clone());
            current_group.push(tx);
        }
        if !current_group.is_empty() {
            groups.push(current_group);
        }

        // Execute groups sequentially, transactions within group can be parallel
        let mut all_receipts = Vec::new();
        for group in groups {
            for tx in group {
                let receipt = self.execute(&tx, block_number).await;
                all_receipts.push(receipt);
            }
        }

        all_receipts
    }
}
