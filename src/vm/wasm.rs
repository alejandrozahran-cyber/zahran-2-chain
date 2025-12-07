use std::collections::HashMap;

pub struct WasmVM {
    instances: HashMap<String, Vec<u8>>,
    gas_limit: u64,
}

impl WasmVM {
    pub fn new(gas_limit: u64) -> Self {
        WasmVM {
            instances: HashMap::new(),
            gas_limit,
        }
    }

    pub fn deploy_contract(&mut self, code: Vec<u8>) -> String {
        let contract_id = format!("wasm_{}", self.instances.len());
        self.instances.insert(contract_id.clone(), code);
        println!("📦 WASM contract deployed: {}", contract_id);
        contract_id
    }

    pub fn execute(&self, contract_id: &str, method: &str, params: Vec<u8>) -> Result<Vec<u8>, String> {
        if ! self.instances.contains_key(contract_id) {
            return Err("Contract not found".to_string());
        }
        
        println!("⚡ Executing WASM: {}::{}", contract_id, method);
        
        // Simulated execution
        Ok(vec![1, 2, 3, 4]) // Mock result
    }

    pub fn get_contract_count(&self) -> usize {
        self.instances.len()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_wasm_vm() {
        let mut vm = WasmVM::new(1_000_000);
        let contract_id = vm.deploy_contract(vec![0x00, 0x61, 0x73, 0x6d]);
        assert_eq!(vm.get_contract_count(), 1);
        
        let result = vm.execute(&contract_id, "test", vec![]);
        assert!(result.is_ok());
    }
}
