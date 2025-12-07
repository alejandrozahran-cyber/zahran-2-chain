"""
NUSA Chain Economic Simulation
Monte Carlo simulation, tokenomics modeling, attack cost analysis
"""

import random
import math
from dataclasses import dataclass
from typing import List

@dataclass
class TokenomicsParams:
    total_supply: int = 1_000_000_000  # 1 billion NUSA
    initial_price: float = 2.5  # $2.50
    inflation_rate: float = 0. 02  # 2% per year
    burn_rate: float = 0. 01  # 1% of tx fees burned
    staking_apy: float = 0.08  # 8% APY
    validator_reward: int = 10  # NUSA per block
    ubi_daily: int = 10  # NUSA per person per day

class EconomicSimulator:
    def __init__(self, params: TokenomicsParams):
        self.params = params
        self.circulating_supply = params.total_supply * 0.3  # 30% initially circulating
        self.price = params.initial_price
        self. total_burned = 0
        self.total_staked = 0
        
    def simulate_year(self, days: int = 365) -> dict:
        """Monte Carlo simulation for 1 year"""
        print(f"📊 Running economic simulation ({days} days)...")
        
        results = {
            "days": [],
            "price": [],
            "supply": [],
            "burned": [],
            "market_cap": [],
        }
        
        for day in range(days):
            # Daily supply changes
            self._simulate_day()
            
            # Record metrics
            results["days"].append(day)
            results["price"]. append(self.price)
            results["supply"].append(self. circulating_supply)
            results["burned"].append(self.total_burned)
            results["market_cap"].append(self.price * self.circulating_supply)
            
            if day % 30 == 0:
                print(f"  Day {day}: Price ${self.price:.2f} | Supply {self.circulating_supply:,. 0f} | MC ${self.price * self.circulating_supply:,.0f}")
        
        return results
    
    def _simulate_day(self):
        """Simulate one day of economic activity"""
        
        # Random daily transactions (normal distribution)
        daily_txs = int(random.gauss(50000, 10000))
        daily_txs = max(10000, daily_txs)  # Minimum 10K txs
        
        # Transaction fees collected
        avg_fee = 0.001  # 0.001 NUSA per tx
        fees_collected = daily_txs * avg_fee
        
        # Burn fees (EIP-1559 style)
        fees_burned = fees_collected * self.params.burn_rate
        self.total_burned += fees_burned
        self.circulating_supply -= fees_burned
        
        # Validator rewards (inflation)
        blocks_per_day = 43200  # 2 sec blocks
        validator_rewards = blocks_per_day * self.params.validator_reward
        self. circulating_supply += validator_rewards
        
        # UBI distribution
        verified_humans = 100000  # Assume 100K verified
        ubi_issued = verified_humans * self.params.ubi_daily
        self.circulating_supply += ubi_issued
        
        # Price movement (random walk with drift)
        daily_volatility = 0.03  # 3% daily volatility
        drift = 0.0003  # Slight upward trend
        price_change = random.gauss(drift, daily_volatility)
        self.price *= (1 + price_change)
        
        # Price floor (can't go below $0.10)
        self.price = max(0.10, self.price)
    
    def supply_demand_analysis(self) -> dict:
        """Analyze supply and demand dynamics"""
        
        # Supply factors
        total_supply = self.params.total_supply
        circulating = self.circulating_supply
        locked = self.total_staked
        burned = self.total_burned
        
        # Demand factors (simplified)
        daily_active_users = 50000
        avg_tx_per_user = 2
        daily_demand = daily_active_users * avg_tx_per_user * 0.001
        
        return {
            "supply": {
                "total": total_supply,
                "circulating": circulating,
                "locked_staking": locked,
                "burned": burned,
                "liquid": circulating - locked,
            },
            "demand": {
                "daily_users": daily_active_users,
                "daily_tx_volume": daily_active_users * avg_tx_per_user,
                "daily_demand_nusa": daily_demand,
            },
            "ratio": {
                "circulating_vs_total": f"{circulating/total_supply*100:.2f}%",
                "burned_vs_total": f"{burned/total_supply*100:.2f}%",
            }
        }
    
    def inflation_curve(self, years: int = 10) -> List[float]:
        """Calculate inflation over time"""
        curve = []
        supply = self.params.total_supply
        
        for year in range(years):
            # Annual inflation from rewards
            annual_blocks = 15768000  # ~365 * 43200
            annual_rewards = annual_blocks * self.params.validator_reward
            
            # Annual burn from fees (assume average)
            annual_burn = annual_rewards * 0.3  # Assume 30% burned
            
            # Net inflation
            net_inflation = (annual_rewards - annual_burn) / supply
            curve.append(net_inflation * 100)  # Percentage
            
            supply += (annual_rewards - annual_burn)
        
        print("📈 Inflation Curve:")
        for i, inflation in enumerate(curve):
            print(f"  Year {i+1}: {inflation:.2f}%")
        
        return curve
    
    def validator_revenue_projection(self, stake_amount: int, days: int = 365) -> dict:
        """Calculate validator earnings"""
        
        # Block rewards
        blocks_per_day = 43200
        blocks_per_year = blocks_per_day * days
        
        # Assume validator gets 1% of blocks (depends on stake)
        validator_share = 0.01
        blocks_produced = int(blocks_per_year * validator_share)
        
        # Rewards
        block_rewards = blocks_produced * self.params.validator_reward
        
        # Transaction fees (validator gets priority fees)
        avg_priority_fee = 0.0005
        avg_txs_per_block = 100
        fee_rewards = blocks_produced * avg_txs_per_block * avg_priority_fee
        
        total_rewards = block_rewards + fee_rewards
        roi = (total_rewards / stake_amount) * 100
        
        return {
            "stake_amount": stake_amount,
            "days": days,
            "blocks_produced": blocks_produced,
            "block_rewards": block_rewards,
            "fee_rewards": fee_rewards,
            "total_rewards": total_rewards,
            "roi": f"{roi:.2f}%",
            "apy": f"{(roi * 365 / days):.2f}%",
        }
    
    def attack_cost_estimation(self) -> dict:
        """Calculate cost of various attacks"""
        
        # 51% attack (need 67% for BFT)
        total_staked = self.params.total_supply * 0.4  # Assume 40% staked
        attack_stake_needed = total_staked * 0.67
        attack_cost_buy = attack_stake_needed * self.price
        
        # Double-spend attack
        avg_block_value = 100000  # $100K per block
        double_spend_profit = avg_block_value
        attack_cost_51 = attack_cost_buy + (avg_block_value * 6)  # 6 blocks deep
        
        # Sybil attack (UBI farming)
        ubi_daily_value = self.params.ubi_daily * self.price
        kyc_cost_per_identity = 100  # $100 to fake identity
        ubi_profit_yearly = ubi_daily_value * 365
        sybil_break_even = kyc_cost_per_identity / ubi_profit_yearly
        
        return {
            "51_percent_attack": {
                "stake_needed": f"{attack_stake_needed:,.0f} NUSA",
                "cost_usd": f"${attack_cost_buy:,.0f}",
                "feasibility": "EXTREMELY DIFFICULT" if attack_cost_buy > 100000000 else "POSSIBLE",
            },
            "double_spend": {
                "potential_profit": f"${double_spend_profit:,.0f}",
                "attack_cost": f"${attack_cost_51:,.0f}",
                "profitable": "NO" if attack_cost_51 > double_spend_profit else "YES",
            },
            "sybil_attack": {
                "ubi_per_day": f"${ubi_daily_value:. 2f}",
                "kyc_cost": f"${kyc_cost_per_identity}",
                "break_even_years": f"{sybil_break_even:.1f}",
                "feasibility": "LOW" if sybil_break_even > 5 else "MODERATE",
            }
        }
    
    def generate_report(self):
        """Generate comprehensive economic report"""
        print("\n" + "="*60)
        print("     NUSA CHAIN ECONOMIC ANALYSIS REPORT")
        print("="*60)
        
        # Supply/demand
        sd = self.supply_demand_analysis()
        print("\n📊 SUPPLY & DEMAND")
        print(f"  Circulating: {sd['supply']['circulating']:,. 0f} NUSA")
        print(f"  Burned: {sd['supply']['burned']:,.0f} NUSA")
        print(f"  Daily Demand: {sd['demand']['daily_demand_nusa']:,.2f} NUSA")
        
        # Inflation
        print("\n📈 INFLATION CURVE (10 years)")
        self.inflation_curve(10)
        
        # Validator
        val = self.validator_revenue_projection(1000000, 365)
        print(f"\n💰 VALIDATOR REVENUE (1M NUSA stake)")
        print(f"  Total Rewards: {val['total_rewards']:,.0f} NUSA")
        print(f"  APY: {val['apy']}")
        
        # Attack costs
        attack = self.attack_cost_estimation()
        print(f"\n🛡️ ATTACK COST ANALYSIS")
        print(f"  51% Attack: {attack['51_percent_attack']['cost_usd']} ({attack['51_percent_attack']['feasibility']})")
        print(f"  Double Spend: {attack['double_spend']['profitable']}")
        print(f"  Sybil Attack: {attack['sybil_attack']['feasibility']} feasibility")
        
        print("\n" + "="*60)

# Run simulation
if __name__ == "__main__":
    params = TokenomicsParams()
    sim = EconomicSimulator(params)
    
    # Run 1 year simulation
    results = sim.simulate_year(365)
    
    print(f"\n📊 FINAL RESULTS:")
    print(f"  Final Price: ${results['price'][-1]:.2f}")
    print(f"  Final Supply: {results['supply'][-1]:,.0f} NUSA")
    print(f"  Total Burned: {results['burned'][-1]:,.0f} NUSA")
    print(f"  Market Cap: ${results['market_cap'][-1]:,. 0f}")
    
    # Generate full report
    sim.generate_report()
