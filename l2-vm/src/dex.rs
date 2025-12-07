// NUSA DEX - Built-in Decentralized Exchange
// Features: AMM, Limit Orders, Zero-Knowledge Swaps

use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct LiquidityPool {
    pub token_a: String,
    pub token_b: String,
    pub reserve_a: u64,
    pub reserve_b: u64,
    pub total_shares: u64,
    pub fee: u64, // 0.3% = 30 basis points
}

impl LiquidityPool {
    pub fn new(token_a: String, token_b: String) -> Self {
        Self {
            token_a,
            token_b,
            reserve_a: 0,
            reserve_b: 0,
            total_shares: 0,
            fee: 30, // 0.3%
        }
    }
    
    // Add liquidity
    pub fn add_liquidity(&mut self, amount_a: u64, amount_b: u64) -> u64 {
        let shares = if self.total_shares == 0 {
            (amount_a * amount_b). sqrt()
        } else {
            std::cmp::min(
                amount_a * self.total_shares / self.reserve_a,
                amount_b * self.total_shares / self.reserve_b,
            )
        };
        
        self.reserve_a += amount_a;
        self. reserve_b += amount_b;
        self.total_shares += shares;
        
        shares
    }
    
    // Swap with constant product formula (x * y = k)
    pub fn swap(&mut self, token_in: String, amount_in: u64) -> u64 {
        let (reserve_in, reserve_out) = if token_in == self.token_a {
            (self.reserve_a, self.reserve_b)
        } else {
            (self.reserve_b, self. reserve_a)
        };
        
        // Apply fee
        let amount_in_with_fee = amount_in * (10000 - self.fee) / 10000;
        
        // Calculate output: y = (y_reserve * x_in) / (x_reserve + x_in)
        let amount_out = (reserve_out * amount_in_with_fee) / (reserve_in + amount_in_with_fee);
        
        // Update reserves
        if token_in == self.token_a {
            self.reserve_a += amount_in;
            self.reserve_b -= amount_out;
        } else {
            self.reserve_b += amount_in;
            self.reserve_a -= amount_out;
        }
        
        amount_out
    }
    
    // Get price
    pub fn get_price(&self, token: String) -> f64 {
        if token == self.token_a {
            self.reserve_b as f64 / self.reserve_a as f64
        } else {
            self.reserve_a as f64 / self.reserve_b as f64
        }
    }
}

pub struct NusaDEX {
    pools: HashMap<String, LiquidityPool>,
    orders: Vec<LimitOrder>,
}

#[derive(Debug, Clone)]
pub struct LimitOrder {
    pub id: String,
    pub trader: String,
    pub token_in: String,
    pub token_out: String,
    pub amount_in: u64,
    pub price: f64,
    pub filled: bool,
}

impl NusaDEX {
    pub fn new() -> Self {
        Self {
            pools: HashMap::new(),
            orders: Vec::new(),
        }
    }
    
    pub fn create_pool(&mut self, token_a: String, token_b: String) {
        let pool_id = format!("{}-{}", token_a, token_b);
        self.pools.insert(pool_id, LiquidityPool::new(token_a, token_b));
    }
    
    pub fn instant_swap(&mut self, pool_id: String, token_in: String, amount: u64) -> u64 {
        if let Some(pool) = self. pools.get_mut(&pool_id) {
            pool.swap(token_in, amount)
        } else {
            0
        }
    }
}
