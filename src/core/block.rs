use serde::{Deserialize, Serialize};
use sha2::{Sha256, Digest};
use chrono::Utc;
use super::transaction::Transaction;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Block {
    pub number: u64,
    pub timestamp: i64,
    pub transactions: Vec<Transaction>,
    pub parent_hash: String,
    pub hash: String,
    pub state_root: String,
    pub gas_used: u64,
    pub gas_limit: u64,
    pub validator: String,
}

impl Block {
    pub fn new(number: u64, parent_hash: String, validator: String) -> Self {
        let mut block = Block {
            number,
            timestamp: Utc::now().timestamp(),
            transactions: vec![],
            parent_hash: parent_hash.clone(),
            hash: String::new(),
            state_root: String::from("0x0"),
            gas_used: 0,
            gas_limit: 30_000_000,
            validator,
        };
        block.hash = block.calculate_hash();
        block
    }

    pub fn add_transaction(&mut self, tx: Transaction) -> bool {
        if self.gas_used + tx.gas_limit <= self.gas_limit {
            self.gas_used += tx.gas_limit;
            self.transactions.push(tx);
            self.hash = self.calculate_hash();
            true
        } else {
            false
        }
    }

    pub fn calculate_hash(&self) -> String {
        let mut hasher = Sha256::new();
        let data = format!(
            "{}{}{}{}",
            self.number,
            self.timestamp,
            self. parent_hash,
            self.transactions.len()
        );
        hasher.update(data.as_bytes());
        format!("{:x}", hasher.finalize())
    }

    pub fn verify(&self) -> bool {
        self.hash == self.calculate_hash()
    }
}

#[derive(Debug, Clone)]
pub struct BlockHeader {
    pub number: u64,
    pub hash: String,
    pub parent_hash: String,
    pub timestamp: i64,
}

impl From<&Block> for BlockHeader {
    fn from(block: &Block) -> Self {
        BlockHeader {
            number: block. number,
            hash: block. hash.clone(),
            parent_hash: block.parent_hash. clone(),
            timestamp: block. timestamp,
        }
    }
}
