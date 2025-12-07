"""NUSA Chain AI Trading Bot"""
import time
from dataclasses import dataclass
from typing import List, Dict
from enum import Enum

class Strategy(Enum):
    MOMENTUM = "momentum"
    MEAN_REVERSION = "mean_reversion"

@dataclass
class Trade:
    timestamp: int
    pair: str
    side: str
    amount: float
    price: float
    profit_loss: float = 0. 0

class AITradingBot:
    def __init__(self, initial_balance: float):
        self. balance = initial_balance
        self.portfolio = {}
        self. trades = []
        self.risk_level = 0.02
        
    def analyze_market(self, price_history: List[float]) -> Dict:
        if len(price_history) < 20:
            return {"signal": "HOLD", "confidence": 0.0}
        
        sma_20 = sum(price_history[-20:]) / 20
        current_price = price_history[-1]
        
        signal = "HOLD"
        confidence = 0.5
        
        if current_price > sma_20:
            signal = "BUY"
            confidence = 0. 8
        elif current_price < sma_20:
            signal = "SELL"
            confidence = 0.8
        
        return {
            "signal": signal,
            "confidence": confidence,
            "current_price": current_price
        }
    
    def execute_trade(self, pair: str, signal: str, amount: float, price: float):
        max_risk = self.balance * self.risk_level
        trade_amount = min(amount, max_risk / price)
        
        if signal == "BUY" and self. balance >= trade_amount * price:
            self.balance -= trade_amount * price
            self.portfolio[pair] = self.portfolio. get(pair, 0) + trade_amount
            
            trade = Trade(
                timestamp=int(time.time()),
                pair=pair,
                side="BUY",
                amount=trade_amount,
                price=price
            )
            self.trades. append(trade)
            return f"✅ BOUGHT {trade_amount:. 4f} {pair}"
        
        return "❌ Trade failed"
    
    def get_performance(self) -> Dict:
        return {
            "total_trades": len(self.trades),
            "balance": self.balance,
            "portfolio": self.portfolio
        }
