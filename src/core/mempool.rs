use std::collections::{HashMap, VecDeque};
use std::sync::Arc;
use tokio::sync::RwLock;
use super::transaction::Transaction;

pub struct Mempool {
    pending: Arc<RwLock<VecDeque<Transaction>>>,
    by_hash: Arc<RwLock<HashMap<String, Transaction>>>,
    max_size: usize,
}

impl Mempool {
    pub fn new(max_size: usize) -> Self {
        Mempool {
            pending: Arc::new(RwLock::new(VecDeque::new())),
            by_hash: Arc::new(RwLock::new(HashMap::new())),
            max_size,
        }
    }

    pub async fn add_transaction(&self, tx: Transaction) -> Result<(), String> {
        let mut pending = self.pending.write().await;
        let mut by_hash = self.by_hash.write().await;

        if pending.len() >= self.max_size {
            return Err("Mempool full".to_string());
        }

        if by_hash.contains_key(&tx.hash) {
            return Err("Transaction already exists".to_string());
        }

        by_hash.insert(tx.hash.clone(), tx. clone());
        pending.push_back(tx);
        
        Ok(())
    }

    pub async fn get_transactions(&self, count: usize) -> Vec<Transaction> {
        let mut pending = self.pending.write().await;
        let mut result = Vec::new();

        for _ in 0..count. min(pending.len()) {
            if let Some(tx) = pending.pop_front() {
                result.push(tx);
            }
        }

        result
    }

    pub async fn size(&self) -> usize {
        self.pending.read().await.len()
    }

    pub async fn remove_transaction(&self, hash: &str) {
        let mut by_hash = self.by_hash.write(). await;
        by_hash. remove(hash);
    }
}
