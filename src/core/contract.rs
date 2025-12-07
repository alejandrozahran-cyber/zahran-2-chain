use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;

pub struct Contract {
    pub address: String,
    pub code: Vec<u8>,
    pub storage: HashMap<String, Vec<u8>>,
}

pub struct ContractEngine {
    contracts: Arc<RwLock<HashMap<String, Contract>>>,
}

impl ContractEngine {
    pub fn new() -> Self {
        ContractEngine {
            contracts: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    pub async fn deploy(&self, code: Vec<u8>) -> String {
        let mut contracts = self.contracts.write().await;
        let address = format!("0xContract{}", contracts.len());
        
        let contract = Contract {
            address: address.clone(),
            code,
            storage: HashMap::new(),
        };
        
        contracts.insert(address. clone(), contract);
        println!("📦 Contract deployed at {}", address);
        
        address
    }

    pub async fn call(&self, address: &str, method: &str, params: Vec<u8>) -> Result<Vec<u8>, String> {
        let contracts = self.contracts.read().await;
        
        if let Some(_contract) = contracts.get(address) {
            println!("⚡ Calling {}::{}", address, method);
            // Simulate execution
            Ok(vec![1, 2, 3, 4])
        } else {
            Err("Contract not found".to_string())
        }
    }

    pub async fn get_code(&self, address: &str) -> Option<Vec<u8>> {
        let contracts = self. contracts.read().await;
        contracts.get(address).map(|c| c.code. clone())
    }
}
