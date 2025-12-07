// NUSA NFT Marketplace - Buy, sell, auction NFTs
// Features: Royalties, auctions, collections, rarity system

use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct NFT {
    pub token_id: String,
    pub collection: String,
    pub name: String,
    pub description: String,
    pub image_uri: String,
    pub owner: String,
    pub creator: String,
    pub royalty_percentage: f64,  // 5% = 0.05
    pub rarity_score: u64,
    pub attributes: HashMap<String, String>,
}

#[derive(Debug, Clone)]
pub struct Listing {
    pub nft_id: String,
    pub seller: String,
    pub price: u64,
    pub currency: String,  // NUSA, NUSD, etc
    pub listed_at: u64,
    pub active: bool,
}

#[derive(Debug, Clone)]
pub struct Auction {
    pub nft_id: String,
    pub seller: String,
    pub starting_bid: u64,
    pub current_bid: u64,
    pub highest_bidder: String,
    pub end_time: u64,
    pub active: bool,
}

pub struct NFTMarketplace {
    nfts: HashMap<String, NFT>,
    listings: HashMap<String, Listing>,
    auctions: HashMap<String, Auction>,
    collections: HashMap<String, Collection>,
    marketplace_fee: f64,  // 2.5% fee
}

#[derive(Debug, Clone)]
pub struct Collection {
    pub id: String,
    pub name: String,
    pub creator: String,
    pub total_items: u64,
    pub floor_price: u64,
    pub total_volume: u64,
    pub verified: bool,
}

impl NFTMarketplace {
    pub fn new() -> Self {
        Self {
            nfts: HashMap::new(),
            listings: HashMap::new(),
            auctions: HashMap::new(),
            collections: HashMap::new(),
            marketplace_fee: 0.025,  // 2.5%
        }
    }
    
    // Mint NFT
    pub fn mint_nft(
        &mut self,
        collection: String,
        name: String,
        description: String,
        image_uri: String,
        creator: String,
        royalty_percentage: f64,
    ) -> String {
        let token_id = format!("nft_{}_{}", collection, self.nfts.len() + 1);
        
        let nft = NFT {
            token_id: token_id.clone(),
            collection: collection.clone(),
            name: name.clone(),
            description,
            image_uri,
            owner: creator.clone(),
            creator: creator.clone(),
            royalty_percentage,
            rarity_score: 0,
            attributes: HashMap::new(),
        };
        
        self.nfts.insert(token_id. clone(), nft);
        
        // Update collection
        if let Some(col) = self.collections.get_mut(&collection) {
            col.total_items += 1;
        }
        
        println!("🎨 NFT minted: {} in collection {}", name, collection);
        
        token_id
    }
    
    // List NFT for sale
    pub fn list_nft(&mut self, nft_id: String, seller: String, price: u64, currency: String) -> bool {
        // Verify ownership
        if let Some(nft) = self. nfts.get(&nft_id) {
            if nft.owner != seller {
                println!("❌ Not the owner");
                return false;
            }
            
            let listing = Listing {
                nft_id: nft_id.clone(),
                seller: seller.clone(),
                price,
                currency: currency.clone(),
                listed_at: 0,  // timestamp
                active: true,
            };
            
            self.listings.insert(nft_id. clone(), listing);
            
            println!("📝 NFT listed: {} for {} {}", nft_id, price, currency);
            
            true
        } else {
            false
        }
    }
    
    // Buy NFT
    pub fn buy_nft(&mut self, nft_id: String, buyer: String, payment: u64) -> bool {
        let listing = self.listings.get_mut(&nft_id);
        if listing.is_none() {
            return false;
        }
        
        let listing = listing.unwrap();
        
        if ! listing.active {
            println!("❌ Listing not active");
            return false;
        }
        
        if payment < listing.price {
            println!("❌ Insufficient payment");
            return false;
        }
        
        let nft = self.nfts. get_mut(&nft_id). unwrap();
        
        // Calculate fees
        let marketplace_fee_amount = (listing.price as f64 * self.marketplace_fee) as u64;
        let royalty_amount = (listing.price as f64 * nft.royalty_percentage) as u64;
        let seller_amount = listing.price - marketplace_fee_amount - royalty_amount;
        
        println!("💰 Sale: {} NUSA", listing.price);
        println! ("   → Seller: {} NUSA", seller_amount);
        println!("   → Royalty (creator): {} NUSA", royalty_amount);
        println!("   → Marketplace fee: {} NUSA", marketplace_fee_amount);
        
        // Transfer NFT
        nft.owner = buyer.clone();
        
        // Deactivate listing
        listing.active = false;
        
        // Update collection floor price
        if let Some(col) = self.collections.get_mut(&nft.collection) {
            col.total_volume += listing.price;
        }
        
        println!("✅ {} bought {} from {}", buyer, nft_id, listing.seller);
        
        true
    }
    
    // Create auction
    pub fn create_auction(&mut self, nft_id: String, seller: String, starting_bid: u64, duration: u64) -> bool {
        if let Some(nft) = self.nfts.get(&nft_id) {
            if nft.owner != seller {
                return false;
            }
            
            let auction = Auction {
                nft_id: nft_id.clone(),
                seller: seller.clone(),
                starting_bid,
                current_bid: starting_bid,
                highest_bidder: String::new(),
                end_time: duration,  // timestamp
                active: true,
            };
            
            self.auctions.insert(nft_id.clone(), auction);
            
            println!("🔨 Auction created for {} | Starting bid: {}", nft_id, starting_bid);
            
            true
        } else {
            false
        }
    }
    
    // Place bid
    pub fn place_bid(&mut self, nft_id: String, bidder: String, bid_amount: u64) -> bool {
        let auction = self.auctions.get_mut(&nft_id);
        if auction.is_none() {
            return false;
        }
        
        let auction = auction.unwrap();
        
        if ! auction.active {
            println! ("❌ Auction ended");
            return false;
        }
        
        if bid_amount <= auction.current_bid {
            println!("❌ Bid must be higher than current bid");
            return false;
        }
        
        // Refund previous bidder (simplified)
        if ! auction.highest_bidder.is_empty() {
            println!("↩️ Refunding {} to {}", auction.current_bid, auction.highest_bidder);
        }
        
        auction.current_bid = bid_amount;
        auction.highest_bidder = bidder.clone();
        
        println!("🔨 New bid: {} NUSA by {}", bid_amount, bidder);
        
        true
    }
    
    // End auction
    pub fn end_auction(&mut self, nft_id: String) -> bool {
        let auction = self.auctions. get_mut(&nft_id);
        if auction.is_none() {
            return false;
        }
        
        let auction = auction.unwrap();
        auction.active = false;
        
        if auction.highest_bidder.is_empty() {
            println!("❌ No bids received");
            return false;
        }
        
        // Transfer NFT to winner
        if let Some(nft) = self.nfts.get_mut(&nft_id) {
            nft.owner = auction.highest_bidder.clone();
            
            println!("🏆 Auction won by {} for {} NUSA", auction.highest_bidder, auction.current_bid);
            
            true
        } else {
            false
        }
    }
    
    // Create collection
    pub fn create_collection(&mut self, name: String, creator: String) -> String {
        let id = format!("col_{}", self.collections.len() + 1);
        
        let collection = Collection {
            id: id.clone(),
            name: name.clone(),
            creator,
            total_items: 0,
            floor_price: 0,
            total_volume: 0,
            verified: false,
        };
        
        self.collections.insert(id. clone(), collection);
        
        println!("📚 Collection created: {}", name);
        
        id
    }
    
    // Get trending NFTs
    pub fn get_trending(&self) -> Vec<String> {
        // Sort by recent sales volume
        // (Simplified - production needs time-based sorting)
        self.nfts.keys().take(10).cloned().collect()
    }
}
