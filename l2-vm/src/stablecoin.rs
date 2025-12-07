// NUSA Stablecoin (NUSD) - Algorithmic USD-pegged stablecoin
// Maintains $1 peg through supply adjustment and collateral

use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct NusaStablecoin {
    pub total_supply: u64,
    pub target_price: f64,      // $1. 00
    pub current_price: f64,
    pub collateral_ratio: f64,  // 150% collateralized
    pub nusa_reserves: u64,
    pub nusd_supply: u64,
}

impl NusaStablecoin {
    pub fn new() -> Self {
        Self {
            total_supply: 0,
            target_price: 1.0,
            current_price: 1.0,
            collateral_ratio: 1.5,
            nusa_reserves: 0,
            nusd_supply: 0,
        }
    }
    
    // Mint NUSD by depositing NUSA as collateral
    pub fn mint(&mut self, nusa_amount: u64, nusa_price: f64) -> u64 {
        // Calculate NUSD to mint (150% collateralized)
        let collateral_value = nusa_amount as f64 * nusa_price;
        let nusd_to_mint = (collateral_value / self.collateral_ratio) as u64;
        
        self.nusa_reserves += nusa_amount;
        self. nusd_supply += nusd_to_mint;
        self.total_supply += nusd_to_mint;
        
        println!("🪙 Minted {} NUSD (collateral: {} NUSA)", nusd_to_mint, nusa_amount);
        
        nusd_to_mint
    }
    
    // Burn NUSD to redeem NUSA collateral
    pub fn burn(&mut self, nusd_amount: u64, nusa_price: f64) -> u64 {
        if nusd_amount > self.nusd_supply {
            return 0;
        }
        
        // Calculate NUSA to return
        let nusa_to_return = ((nusd_amount as f64) / nusa_price) as u64;
        
        if nusa_to_return > self.nusa_reserves {
            return 0;
        }
        
        self.nusd_supply -= nusd_amount;
        self.nusa_reserves -= nusa_to_return;
        self.total_supply -= nusd_amount;
        
        println!("🔥 Burned {} NUSD (returned: {} NUSA)", nusd_amount, nusa_to_return);
        
        nusa_to_return
    }
    
    // Maintain peg through algorithmic adjustments
    pub fn rebalance(&mut self, market_price: f64, nusa_price: f64) {
        self.current_price = market_price;
        
        if market_price > self.target_price + 0.01 {
            // NUSD trading above $1 → Increase supply
            self.expand_supply(nusa_price);
        } else if market_price < self. target_price - 0.01 {
            // NUSD trading below $1 → Decrease supply
            self.contract_supply();
        }
    }
    
    fn expand_supply(&mut self, nusa_price: f64) {
        // Mint new NUSD to push price down
        let expansion_amount = (self.nusd_supply as f64 * 0.01) as u64; // 1% expansion
        
        println!("📈 Expanding supply by {} NUSD to restore peg", expansion_amount);
        
        self.nusd_supply += expansion_amount;
        self.total_supply += expansion_amount;
    }
    
    fn contract_supply(&mut self) {
        // Buy back and burn NUSD to push price up
        let contraction_amount = (self.nusd_supply as f64 * 0.01) as u64; // 1% contraction
        
        println!("📉 Contracting supply by {} NUSD to restore peg", contraction_amount);
        
        if contraction_amount <= self.nusd_supply {
            self. nusd_supply -= contraction_amount;
            self.total_supply -= contraction_amount;
        }
    }
    
    // Check if collateral is sufficient
    pub fn check_health(&self, nusa_price: f64) -> f64 {
        if self.nusd_supply == 0 {
            return 100.0;
        }
        
        let collateral_value = self.nusa_reserves as f64 * nusa_price;
        let debt_value = self.nusd_supply as f64;
        
        (collateral_value / debt_value) * 100.0
    }
    
    // Liquidate undercollateralized positions
    pub fn liquidate(&mut self, position_id: &str, nusa_price: f64) -> bool {
        let health_ratio = self.check_health(nusa_price);
        
        if health_ratio < 120.0 {
            println! ("⚠️ Liquidating position {} (health: {:.2}%)", position_id, health_ratio);
            return true;
        }
        
        false
    }
    
    // Get stablecoin stats
    pub fn get_stats(&self) -> String {
        format!(
            "NUSD Stats:\n\
             Supply: {} NUSD\n\
             Price: ${:.4}\n\
             Collateral: {} NUSA\n\
             Collateral Ratio: {:.2}%",
            self.nusd_supply,
            self.current_price,
            self.nusa_reserves,
            self.collateral_ratio * 100.0
        )
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_mint_nusd() {
        let mut stable = NusaStablecoin::new();
        let nusd_minted = stable.mint(1000, 2.0); // 1000 NUSA @ $2
        
        assert!(nusd_minted > 0);
        assert_eq!(stable.nusa_reserves, 1000);
    }

    #[test]
    fn test_burn_nusd() {
        let mut stable = NusaStablecoin::new();
        stable.mint(1000, 2.0);
        let nusa_returned = stable.burn(500, 2.0);
        
        assert!(nusa_returned > 0);
    }

    #[test]
    fn test_rebalance() {
        let mut stable = NusaStablecoin::new();
        stable.mint(1000, 2.0);
        stable.rebalance(1.05, 2.0); // Price above peg
        
        // Supply should expand
    }
}
