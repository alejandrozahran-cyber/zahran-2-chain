use serde::{Deserialize, Serialize};
use sha2::{Sha256, Digest};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Transaction {
    pub from: String,
    pub to: String,
    pub value: u64,
    pub gas_price: u64,
    pub gas_limit: u64,
    pub nonce: u64,
    pub data: Vec<u8>,
    pub signature: Vec<u8>,
    pub hash: String,
}

impl Transaction {
    pub fn new(from: String, to: String, value: u64, nonce: u64) -> Self {
        let mut tx = Transaction {
            from: from. clone(),
            to: to. clone(),
            value,
            gas_price: 1_000_000_000, // 1 Gwei
            gas_limit: 21000,
            nonce,
            data: vec![],
            signature: vec![],
            hash: String::new(),
        };
        tx.hash = tx.calculate_hash();
        tx
    }

    pub fn calculate_hash(&self) -> String {
        let mut hasher = Sha256::new();
        let data = format!("{}{}{}{}", self.from, self.to, self.value, self.nonce);
        hasher.update(data.as_bytes());
        format!("{:x}", hasher. finalize())
    }

    pub fn verify(&self) -> bool {
        // Simplified verification
        ! self.from.is_empty() && !self.to.is_empty() && self.value > 0
    }
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct TransactionReceipt {
    pub tx_hash: String,
    pub block_number: u64,
    pub gas_used: u64,
    pub status: bool,
    pub logs: Vec<String>,
}
