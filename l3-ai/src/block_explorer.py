"""
NUSA Chain Block Explorer
View blocks, transactions, addresses, contracts, stats in real-time
"""

from dataclasses import dataclass
from typing import List, Dict, Optional
from datetime import datetime
import time

@dataclass
class Block:
    number: int
    hash: str
    parent_hash: str
    timestamp: int
    miner: str
    transactions: List[str]
    gas_used: int
    gas_limit: int
    difficulty: int
    size: int
    reward: int

@dataclass
class Transaction:
    hash: str
    from_addr: str
    to_addr: str
    value: int
    gas_price: int
    gas_used: int
    nonce: int
    block_number: int
    timestamp: int
    status: str  # success, failed, pending
    input_data: str

@dataclass
class Address:
    address: str
    balance: int
    tx_count: int
    first_seen: int
    last_seen: int
    is_contract: bool
    contract_code: Optional[str]

class BlockExplorer:
    def __init__(self):
        self.blocks = {}
        self.transactions = {}
        self.addresses = {}
        self.contracts = {}
        self.total_supply = 1_000_000_000 * 10**8
        self.circulating_supply = 0
        
    # === BLOCK QUERIES ===
    
    def get_block(self, block_number: int) -> Optional[Dict]:
        """Get block by number"""
        block = self.blocks.get(block_number)
        if not block:
            return None
            
        return {
            "number": block. number,
            "hash": block.hash,
            "parent_hash": block.parent_hash,
            "timestamp": datetime.fromtimestamp(block. timestamp).isoformat(),
            "miner": block.miner,
            "transactions": len(block.transactions),
            "gas_used": block.gas_used,
            "gas_limit": block.gas_limit,
            "utilization": f"{block.gas_used / block.gas_limit * 100:. 2f}%",
            "difficulty": block.difficulty,
            "size": f"{block.size / 1024:.2f} KB",
            "reward": block.reward,
            "tx_list": block.transactions[:10]  # First 10 txs
        }
    
    def get_latest_blocks(self, count: int = 10) -> List[Dict]:
        """Get latest N blocks"""
        latest = sorted(self.blocks.values(), key=lambda b: b.number, reverse=True)[:count]
        
        return [{
            "number": block.number,
            "hash": block. hash[:16] + ".. .",
            "timestamp": datetime.fromtimestamp(block.timestamp).isoformat(),
            "miner": block. miner[:16] + "...",
            "txs": len(block.transactions),
            "gas_used": block.gas_used,
            "size": f"{block.size / 1024:.2f} KB"
        } for block in latest]
    
    # === TRANSACTION QUERIES ===
    
    def get_transaction(self, tx_hash: str) -> Optional[Dict]:
        """Get transaction by hash"""
        tx = self.transactions.get(tx_hash)
        if not tx:
            return None
            
        return {
            "hash": tx.hash,
            "status": tx.status,
            "block": tx.block_number,
            "timestamp": datetime.fromtimestamp(tx.timestamp).isoformat(),
            "from": tx.from_addr,
            "to": tx.to_addr,
            "value": f"{tx.value / 10**8:.8f} NUSA",
            "gas_price": tx.gas_price,
            "gas_used": tx. gas_used,
            "tx_fee": f"{tx.gas_price * tx.gas_used / 10**8:.8f} NUSA",
            "nonce": tx.nonce,
            "input_data": tx.input_data[:100] + "..." if len(tx. input_data) > 100 else tx.input_data
        }
    
    def get_pending_transactions(self) -> List[Dict]:
        """Get pending transactions from mempool"""
        pending = [tx for tx in self.transactions. values() if tx.status == "pending"]
        
        return [{
            "hash": tx.hash[:16] + "...",
            "from": tx.from_addr[:16] + "...",
            "to": tx.to_addr[:16] + "...",
            "value": f"{tx. value / 10**8:.4f} NUSA",
            "gas_price": tx.gas_price,
            "age": f"{(time.time() - tx.timestamp):.0f}s"
        } for tx in sorted(pending, key=lambda t: t.gas_price, reverse=True)[:20]]
    
    # === ADDRESS QUERIES ===
    
    def get_address(self, address: str) -> Optional[Dict]:
        """Get address details"""
        addr = self.addresses.get(address)
        if not addr:
            return None
            
        return {
            "address": addr.address,
            "balance": f"{addr.balance / 10**8:.8f} NUSA",
            "usd_value": f"${addr.balance / 10**8 * 2. 5:.2f}",  # Assuming $2.5/NUSA
            "transactions": addr.tx_count,
            "first_seen": datetime.fromtimestamp(addr. first_seen).isoformat(),
            "last_seen": datetime.fromtimestamp(addr. last_seen).isoformat(),
            "is_contract": addr.is_contract,
            "type": "Contract" if addr.is_contract else "Wallet"
        }
    
    def get_address_transactions(self, address: str, page: int = 1, limit: int = 20) -> List[Dict]:
        """Get transactions for an address"""
        txs = [tx for tx in self.transactions. values() 
               if tx.from_addr == address or tx.to_addr == address]
        
        txs_sorted = sorted(txs, key=lambda t: t.timestamp, reverse=True)
        start = (page - 1) * limit
        end = start + limit
        
        return [{
            "hash": tx.hash[:16] + "...",
            "block": tx.block_number,
            "timestamp": datetime.fromtimestamp(tx.timestamp).isoformat(),
            "from": tx.from_addr[:16] + "...",
            "to": tx.to_addr[:16] + "...",
            "value": f"{tx.value / 10**8:.4f} NUSA",
            "direction": "IN" if tx.to_addr == address else "OUT",
            "status": tx.status
        } for tx in txs_sorted[start:end]]
    
    def get_rich_list(self, limit: int = 100) -> List[Dict]:
        """Get richest addresses"""
        sorted_addrs = sorted(self.addresses. values(), key=lambda a: a.balance, reverse=True)[:limit]
        
        total_balance = sum(a.balance for a in sorted_addrs)
        
        return [{
            "rank": i + 1,
            "address": addr.address[:16] + ".. .",
            "balance": f"{addr.balance / 10**8:. 2f} NUSA",
            "percentage": f"{addr.balance / self.circulating_supply * 100:.4f}%",
            "type": "Contract" if addr.is_contract else "Wallet"
        } for i, addr in enumerate(sorted_addrs)]
    
    # === STATS & ANALYTICS ===
    
    def get_network_stats(self) -> Dict:
        """Get overall network statistics"""
        latest_block = max(self.blocks.keys()) if self.blocks else 0
        total_txs = len(self.transactions)
        pending_txs = len([tx for tx in self.transactions.values() if tx.status == "pending"])
        
        # Calculate average block time
        if len(self.blocks) > 1:
            recent_blocks = sorted(self.blocks.values(), key=lambda b: b.number, reverse=True)[:100]
            if len(recent_blocks) > 1:
                time_diff = recent_blocks[0].timestamp - recent_blocks[-1].timestamp
                avg_block_time = time_diff / (len(recent_blocks) - 1)
            else:
                avg_block_time = 2. 0
        else:
            avg_block_time = 2.0
        
        # Calculate TPS
        if len(self.blocks) > 10:
            recent = sorted(self.blocks.values(), key=lambda b: b.number, reverse=True)[:10]
            total_txs_recent = sum(len(b.transactions) for b in recent)
            time_span = recent[0].timestamp - recent[-1].timestamp
            tps = total_txs_recent / time_span if time_span > 0 else 0
        else:
            tps = 0
        
        return {
            "latest_block": latest_block,
            "total_transactions": total_txs,
            "pending_transactions": pending_txs,
            "total_addresses": len(self.addresses),
            "total_contracts": len(self.contracts),
            "avg_block_time": f"{avg_block_time:.2f}s",
            "tps": f"{tps:.2f}",
            "total_supply": f"{self.total_supply / 10**8:,.0f} NUSA",
            "circulating_supply": f"{self.circulating_supply / 10**8:,.0f} NUSA",
            "market_cap": f"${self. circulating_supply / 10**8 * 2.5:,.2f}",
            "price": "$2.50"
        }
    
    def get_gas_stats(self) -> Dict:
        """Get gas price statistics"""
        recent_txs = sorted(self.transactions. values(), 
                           key=lambda t: t.timestamp, reverse=True)[:1000]
        
        if not recent_txs:
            return {"error": "No transactions"}
        
        gas_prices = [tx.gas_price for tx in recent_txs if tx.status != "pending"]
        
        if not gas_prices:
            return {"error": "No completed transactions"}
        
        gas_prices.sort()
        
        return {
            "slow": gas_prices[len(gas_prices) // 10],  # 10th percentile
            "standard": gas_prices[len(gas_prices) // 2],  # 50th percentile (median)
            "fast": gas_prices[len(gas_prices) * 9 // 10],  # 90th percentile
            "instant": gas_prices[-1],  # Maximum
            "avg": sum(gas_prices) // len(gas_prices)
        }
    
    def search(self, query: str) -> Dict:
        """Universal search: block number, tx hash, or address"""
        query = query.strip()
        
        # Try as block number
        if query.isdigit():
            block_num = int(query)
            if block_num in self.blocks:
                return {"type": "block", "data": self.get_block(block_num)}
        
        # Try as transaction hash
        if query in self.transactions:
            return {"type": "transaction", "data": self.get_transaction(query)}
        
        # Try as address
        if query in self.addresses:
            return {"type": "address", "data": self.get_address(query)}
        
        return {"type": "not_found", "message": f"No results for: {query}"}
    
    # === DATA INSERTION (for testing) ===
    
    def add_block(self, block: Block):
        """Add block to explorer"""
        self.blocks[block.number] = block
        
    def add_transaction(self, tx: Transaction):
        """Add transaction to explorer"""
        self.transactions[tx.hash] = tx
        
        # Update addresses
        if tx. from_addr not in self.addresses:
            self.addresses[tx.from_addr] = Address(
                address=tx.from_addr,
                balance=0,
                tx_count=0,
                first_seen=tx.timestamp,
                last_seen=tx.timestamp,
                is_contract=False,
                contract_code=None
            )
        
        if tx.to_addr not in self.addresses:
            self.addresses[tx.to_addr] = Address(
                address=tx.to_addr,
                balance=0,
                tx_count=0,
                first_seen=tx.timestamp,
                last_seen=tx. timestamp,
                is_contract=False,
                contract_code=None
            )
        
        # Update tx counts
        self.addresses[tx.from_addr].tx_count += 1
        self.addresses[tx.to_addr].tx_count += 1
        
        # Update last seen
        self.addresses[tx.from_addr].last_seen = tx.timestamp
        self.addresses[tx.to_addr].last_seen = tx.timestamp


# Example usage
if __name__ == "__main__":
    explorer = BlockExplorer()
    
    # Add test data
    test_block = Block(
        number=12345,
        hash="0xabc123...",
        parent_hash="0xdef456...",
        timestamp=int(time.time()),
        miner="nusa1validator1",
        transactions=["tx1", "tx2", "tx3"],
        gas_used=15000000,
        gas_limit=30000000,
        difficulty=1000,
        size=50000,
        reward=10 * 10**8
    )
    
    explorer.add_block(test_block)
    
    # Query
    print("📊 Network Stats:")
    print(explorer.get_network_stats())
    
    print("\n🔍 Latest Blocks:")
    print(explorer.get_latest_blocks(5))
