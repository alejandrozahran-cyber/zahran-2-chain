use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;

pub struct ERC20Token {
    name: String,
    symbol: String,
    total_supply: u64,
    balances: Arc<RwLock<HashMap<String, u64>>>,
}

impl ERC20Token {
    pub fn new(name: String, symbol: String, total_supply: u64) -> Self {
        let mut balances = HashMap::new();
        balances.insert("0x0000000000000000000000000000000000000000".to_string(), total_supply);
        
        ERC20Token {
            name,
            symbol,
            total_supply,
            balances: Arc::new(RwLock::new(balances)),
        }
    }

    pub async fn balance_of(&self, address: &str) -> u64 {
        let balances = self.balances. read().await;
        *balances.get(address).unwrap_or(&0)
    }

    pub async fn transfer(&self, from: &str, to: &str, amount: u64) -> Result<bool, String> {
        let mut balances = self.balances. write().await;
        
        let from_balance = *balances.get(from).unwrap_or(&0);
        if from_balance < amount {
            return Err("Insufficient balance". to_string());
        }
        
        balances.insert(from.to_string(), from_balance - amount);
        let to_balance = *balances. get(to).unwrap_or(&0);
        balances.insert(to.to_string(), to_balance + amount);
        
        Ok(true)
    }

    pub fn name(&self) -> &str {
        &self.name
    }

    pub fn symbol(&self) -> &str {
        &self.symbol
    }

    pub fn total_supply(&self) -> u64 {
        self.total_supply
    }
}

pub struct ERC721Token {
    name: String,
    symbol: String,
    owners: Arc<RwLock<HashMap<u64, String>>>,
}

impl ERC721Token {
    pub fn new(name: String, symbol: String) -> Self {
        ERC721Token {
            name,
            symbol,
            owners: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    pub async fn mint(&self, token_id: u64, owner: String) -> Result<bool, String> {
        let mut owners = self.owners.write(). await;
        
        if owners.contains_key(&token_id) {
            return Err("Token already exists".to_string());
        }
        
        owners.insert(token_id, owner);
        Ok(true)
    }

    pub async fn owner_of(&self, token_id: u64) -> Option<String> {
        let owners = self.owners.read().await;
        owners.get(&token_id).cloned()
    }

    pub async fn transfer(&self, token_id: u64, to: String) -> Result<bool, String> {
        let mut owners = self. owners.write().await;
        
        if !owners.contains_key(&token_id) {
            return Err("Token doesn't exist".to_string());
        }
        
        owners.insert(token_id, to);
        Ok(true)
    }
}
