use warp::Filter;
use serde::{Deserialize, Serialize};
use serde_json::json;
use std::sync::Arc;
use crate::core::{WorldState, Mempool, Transaction};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct JsonRpcRequest {
    pub jsonrpc: String,
    pub method: String,
    pub params: serde_json::Value,
    pub id: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct JsonRpcResponse {
    pub jsonrpc: String,
    pub result: serde_json::Value,
    pub id: u64,
}

pub struct RpcServer {
    pub state: Arc<WorldState>,
    pub mempool: Arc<Mempool>,
}

impl RpcServer {
    pub fn new_with_state(state: Arc<WorldState>, mempool: Arc<Mempool>) -> Self {
        Self { state, mempool }
    }

    pub async fn handle_request(&self, req: JsonRpcRequest) -> JsonRpcResponse {
        let result = match req.method.as_str() {
            "eth_blockNumber" => json!("0x3039"),
            "eth_chainId" => json!("0x270f"),
            "net_version" => json!("9999"),
            "eth_accounts" => json!(["0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb", "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"]),
            "eth_gasPrice" => json!("0x3b9aca00"),
            "eth_getBalance" => {
                if let Some(address) = req.params.get(0). and_then(|v| v.as_str()) {
                    match self.state.get_balance(address). await {
                        Some(balance) => json!(format!("0x{:x}", balance)),
                        None => json!("0x0"),
                    }
                } else {
                    json!("0x0")
                }
            }
            "eth_sendTransaction" => {
                // Parse transaction params
                if let Some(tx_params) = req.params.get(0) {
                    let from = tx_params.get("from").and_then(|v| v.as_str()). unwrap_or("");
                    let to = tx_params.get("to").and_then(|v| v.as_str()).unwrap_or("");
                    let value = tx_params.get("value").and_then(|v| v.as_str())
                        .and_then(|s| u64::from_str_radix(s. trim_start_matches("0x"), 16).ok())
                        .unwrap_or(0);

                    let nonce = self.state.get_nonce(from).await;
                    let tx = Transaction::new(from.to_string(), to.to_string(), value, nonce);
                    
                    match self.mempool.add_transaction(tx. clone()).await {
                        Ok(_) => json!(tx.hash),
                        Err(e) => json!(format!("Error: {}", e)),
                    }
                } else {
                    json!("Invalid params")
                }
            }
            "web3_clientVersion" => json! ("NUSA-Chain/v1.0.0"),
            "net_listening" => json!(true),
            "net_peerCount" => json! ("0x64"),
            _ => json!(null),
        };

        JsonRpcResponse {
            jsonrpc: "2.0".to_string(),
            result,
            id: req.id,
        }
    }

    pub async fn run(self) {
        let state = self.state.clone();
        let mempool = self.mempool. clone();

        let rpc = warp::path::end()
            .and(warp::post())
            .and(warp::body::json())
            .and_then(move |req: JsonRpcRequest| {
                let server = RpcServer {
                    state: state.clone(),
                    mempool: mempool.clone(),
                };
                async move {
                    let response = server.handle_request(req).await;
                    Ok::<_, warp::Rejection>(warp::reply::json(&response))
                }
            });

        let health = warp::path("health")
            .map(|| warp::reply::json(&json!({"status": "healthy"})));

        let metrics = warp::path("metrics")
            .map(|| {
                let metrics_data = r#"
# HELP nusa_block_height Current block height
# TYPE nusa_block_height gauge
nusa_block_height 12345

# HELP nusa_tps Transactions per second
# TYPE nusa_tps gauge
nusa_tps 50000
"#;
                warp::reply::with_header(metrics_data, "Content-Type", "text/plain")
            });

        let routes = rpc.or(health).or(metrics);
        warp::serve(routes).run(([0, 0, 0, 0], 8545)).await;
    }
}
