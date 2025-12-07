// NUSA Lending Protocol - Borrow, Lend, Earn Interest
// Features: Collateralized loans, flash loans, liquidations

use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct LendingPool {
    pub asset: String,
    pub total_supplied: u64,
    pub total_borrowed: u64,
    pub supply_apy: f64,
    pub borrow_apy: f64,
    pub utilization_rate: f64,
    pub collateral_factor: f64,  // 75% = can borrow up to 75% of collateral
}

#[derive(Debug, Clone)]
pub struct UserPosition {
    pub user: String,
    pub supplied: HashMap<String, u64>,
    pub borrowed: HashMap<String, u64>,
    pub collateral: HashMap<String, u64>,
    pub health_factor: f64,  // Must be > 1.0 to avoid liquidation
}

pub struct LendingProtocol {
    pools: HashMap<String, LendingPool>,
    positions: HashMap<String, UserPosition>,
    liquidation_threshold: f64,  // 1.2 = 120%
    liquidation_bonus: f64,      // 5% bonus for liquidators
}

impl LendingProtocol {
    pub fn new() -> Self {
        let mut protocol = Self {
            pools: HashMap::new(),
            positions: HashMap::new(),
            liquidation_threshold: 1.2,
            liquidation_bonus: 0.05,
        };
        
        // Initialize default pools
        protocol.create_pool("NUSA".to_string(), 0.75);
        protocol.create_pool("NUSD".to_string(), 0.80);
        protocol.create_pool("ETH".to_string(), 0.70);
        
        protocol
    }
    
    // Create lending pool
    pub fn create_pool(&mut self, asset: String, collateral_factor: f64) {
        let pool = LendingPool {
            asset: asset.clone(),
            total_supplied: 0,
            total_borrowed: 0,
            supply_apy: 3.0,   // 3% APY for suppliers
            borrow_apy: 8.0,   // 8% APY for borrowers
            utilization_rate: 0.0,
            collateral_factor,
        };
        
        self.pools.insert(asset. clone(), pool);
        println! ("🏦 Lending pool created: {}", asset);
    }
    
    // Supply assets to earn interest
    pub fn supply(&mut self, user: String, asset: String, amount: u64) -> bool {
        if let Some(pool) = self. pools.get_mut(&asset) {
            pool.total_supplied += amount;
            
            // Update user position
            let position = self.positions.entry(user.clone()).or_insert(UserPosition {
                user: user.clone(),
                supplied: HashMap::new(),
                borrowed: HashMap::new(),
                collateral: HashMap::new(),
                health_factor: 100.0,
            });
            
            let current = position.supplied.get(&asset).unwrap_or(&0);
            position.supplied.insert(asset.clone(), current + amount);
            
            // Update utilization rate
            self.update_pool_rates(&asset);
            
            println!("💵 {} supplied {} {} | APY: {:.2}%", user, amount, asset, pool.supply_apy);
            
            true
        } else {
            false
        }
    }
    
    // Borrow assets (must have collateral)
    pub fn borrow(&mut self, user: String, asset: String, amount: u64) -> bool {
        // Check if pool has liquidity
        if let Some(pool) = self.pools.get_mut(&asset) {
            let available = pool.total_supplied - pool.total_borrowed;
            if available < amount {
                println!("❌ Insufficient liquidity in pool");
                return false;
            }
            
            // Check user's borrowing power
            let position = self.positions.get_mut(&user);
            if position.is_none() {
                println! ("❌ User has no collateral");
                return false;
            }
            
            let position = position.unwrap();
            let borrow_power = self.calculate_borrow_power(&position);
            let current_borrowed = self.calculate_total_borrowed(&position);
            
            if current_borrowed + (amount as f64) > borrow_power {
                println!("❌ Insufficient collateral to borrow");
                return false;
            }
            
            // Execute borrow
            pool.total_borrowed += amount;
            let current = position.borrowed.get(&asset). unwrap_or(&0);
            position.borrowed.insert(asset.clone(), current + amount);
            
            // Update health factor
            position.health_factor = self.calculate_health_factor(&position);
            
            // Update rates
            self.update_pool_rates(&asset);
            
            println!("💳 {} borrowed {} {} | APY: {:.2}% | Health: {:.2}", 
                user, amount, asset, pool.borrow_apy, position.health_factor);
            
            true
        } else {
            false
        }
    }
    
    // Deposit collateral
    pub fn deposit_collateral(&mut self, user: String, asset: String, amount: u64) {
        let position = self.positions.entry(user.clone()).or_insert(UserPosition {
            user: user.clone(),
            supplied: HashMap::new(),
            borrowed: HashMap::new(),
            collateral: HashMap::new(),
            health_factor: 100.0,
        });
        
        let current = position.collateral.get(&asset).unwrap_or(&0);
        position.collateral.insert(asset.clone(), current + amount);
        
        position.health_factor = self.calculate_health_factor(&position);
        
        println!("🔒 {} deposited {} {} as collateral", user, amount, asset);
    }
    
    // Repay borrowed assets
    pub fn repay(&mut self, user: String, asset: String, amount: u64) -> bool {
        if let Some(pool) = self.pools. get_mut(&asset) {
            let position = self.positions.get_mut(&user);
            if position. is_none() {
                return false;
            }
            
            let position = position.unwrap();
            let borrowed = position.borrowed.get(&asset).unwrap_or(&0);
            
            let repay_amount = std::cmp::min(amount, *borrowed);
            
            pool.total_borrowed -= repay_amount;
            position. borrowed.insert(asset.clone(), borrowed - repay_amount);
            
            position.health_factor = self.calculate_health_factor(&position);
            
            self.update_pool_rates(&asset);
            
            println!("✅ {} repaid {} {} | Health: {:.2}", user, repay_amount, asset, position. health_factor);
            
            true
        } else {
            false
        }
    }
    
    // Flash loan (borrow & repay in same transaction)
    pub fn flash_loan(&mut self, asset: String, amount: u64) -> Result<(), String> {
        if let Some(pool) = self. pools.get(&asset) {
            let available = pool.total_supplied - pool.total_borrowed;
            if available < amount {
                return Err("Insufficient liquidity".to_string());
            }
            
            // Flash loan fee: 0.09%
            let fee = (amount as f64 * 0.0009) as u64;
            
            println!("⚡ Flash loan: {} {} (fee: {})", amount, asset, fee);
            
            // User must repay + fee in same transaction
            // (Production: Execute user's arbitrage logic here)
            
            Ok(())
        } else {
            Err("Pool not found".to_string())
        }
    }
    
    // Liquidate undercollateralized position
    pub fn liquidate(&mut self, liquidator: String, user: String, asset: String) -> bool {
        let position = self.positions.get_mut(&user);
        if position.is_none() {
            return false;
        }
        
        let position = position.unwrap();
        
        // Check if liquidatable (health factor < 1.2)
        if position.health_factor >= self.liquidation_threshold {
            println!("❌ Position is healthy, cannot liquidate");
            return false;
        }
        
        let borrowed = position.borrowed.get(&asset).unwrap_or(&0);
        if *borrowed == 0 {
            return false;
        }
        
        // Liquidator pays debt, gets collateral + bonus
        let liquidation_amount = *borrowed;
        let bonus = (liquidation_amount as f64 * self.liquidation_bonus) as u64;
        
        println!("⚠️ LIQUIDATION: {} liquidating {} | Debt: {} | Bonus: {}", 
            liquidator, user, liquidation_amount, bonus);
        
        // Clear debt
        position.borrowed.insert(asset.clone(), 0);
        
        // Transfer collateral to liquidator
        // (Simplified - production needs full collateral management)
        
        position.health_factor = self.calculate_health_factor(&position);
        
        true
    }
    
    // Calculate borrowing power based on collateral
    fn calculate_borrow_power(&self, position: &UserPosition) -> f64 {
        let mut total_collateral_value = 0.0;
        
        for (asset, amount) in &position.collateral {
            if let Some(pool) = self.pools.get(asset) {
                // Simplified: Assume 1:1 price (production needs oracle)
                let value = *amount as f64;
                total_collateral_value += value * pool.collateral_factor;
            }
        }
        
        total_collateral_value
    }
    
    // Calculate total borrowed value
    fn calculate_total_borrowed(&self, position: &UserPosition) -> f64 {
        let mut total = 0.0;
        for (_, amount) in &position.borrowed {
            total += *amount as f64;
        }
        total
    }
    
    // Calculate health factor
    fn calculate_health_factor(&self, position: &UserPosition) -> f64 {
        let total_borrowed = self.calculate_total_borrowed(position);
        if total_borrowed == 0.0 {
            return 100.0;
        }
        
        let borrow_power = self.calculate_borrow_power(position);
        borrow_power / total_borrowed
    }
    
    // Update pool interest rates based on utilization
    fn update_pool_rates(&mut self, asset: &str) {
        if let Some(pool) = self. pools.get_mut(asset) {
            if pool.total_supplied == 0 {
                return;
            }
            
            // Utilization rate = borrowed / supplied
            pool.utilization_rate = pool.total_borrowed as f64 / pool.total_supplied as f64;
            
            // Dynamic interest rates
            pool.borrow_apy = 2.0 + (pool.utilization_rate * 20.0); // 2-22% APY
            pool.supply_apy = pool.b*

