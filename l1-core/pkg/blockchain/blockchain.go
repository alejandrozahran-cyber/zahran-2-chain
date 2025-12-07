package blockchain

import (
	"errors"
	"sync"
)

var (
	ErrInvalidBlock = errors.New("invalid block")
)

type Blockchain struct {
	Blocks   []*Block
	Balances map[string]float64
	mu       sync.RWMutex
}

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Blocks:   make([]*Block, 0),
		Balances: make(map[string]float64),
	}
	bc. Blocks = append(bc.Blocks, CreateGenesisBlock())
	bc.Balances["genesis"] = 25000000. 0
	return bc
}

func (bc *Blockchain) GetLatestBlock() *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	if len(bc.Blocks) == 0 {
		return nil
	}
	return bc. Blocks[len(bc. Blocks)-1]
}

func (bc *Blockchain) GetHeight() uint64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return uint64(len(bc.Blocks))
}

func (bc *Blockchain) GetBalance(address string) float64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc. Balances[address]
}

func (bc *Blockchain) IsValid() bool {
	return len(bc.Blocks) > 0
}

func (bc *Blockchain) GetAllBlocks() []*Block {
	bc. mu.RLock()
	defer bc.mu.RUnlock()
	blocks := make([]*Block, len(bc.Blocks))
	copy(blocks, bc.Blocks)
	return blocks
}
