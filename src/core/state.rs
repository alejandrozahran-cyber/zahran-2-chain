use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;

#[derive(Debug, Clone)]
pub struct Account {
    pub address: String,
    pub balance: u64,
    pub nonce: u64,
    pub code: Vec<u8>,
    pub storage: HashMap<String, String>,
}

impl Account {
    pub fn new(address: String, balance: u64) -> Self {
        Account {
            address,
            balance,
            nonce: 0,
            code: vec![],
            storage: HashMap::new(),
        }
    }
}

pub struct WorldState {
    accounts: Arc<RwLock<HashMap<String, Account>>>,
}

impl WorldState {
    pub fn new() -> Self {
        WorldState {
            accounts: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    pub async fn create_account(&self, address: String, balance: u64) {
        let mut accounts = self. accounts.write().await;
        accounts.insert(address. clone(), Account::new(address, balance));
    }

    pub async fn get_balance(&self, address: &str) -> Option<u64> {
        let accounts = self.accounts. read().await;
        accounts. get(address).map(|acc| acc.balance)
    }

    pub async fn transfer(&self, from: &str, to: &str, amount: u64) -> Result<(), String> {
        let mut accounts = self.accounts.write().await;
        
        let from_balance = accounts.get(from).ok_or("From account not found")?.balance;
        if from_balance < amount {
            return Err("Insufficient balance".to_string());
        }

        // Deduct from sender
        if let Some(from_acc) = accounts.get_mut(from) {
            from_acc.balance -= amount;
            from_acc.nonce += 1;
        }

        // Add to receiver (create if doesn't exist)
        accounts.entry(to.to_string())
            .and_modify(|acc| acc.balance += amount)
            .or_insert_with(|| Account::new(to.to_string(), amount));

        Ok(())
    }

    pub async fn get_nonce(&self, address: &str) -> u64 {
        let accounts = self.accounts.read().await;
        accounts.get(address). map(|acc| acc.nonce).unwrap_or(0)
    }
}
