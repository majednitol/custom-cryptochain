package blockchain

import (
	"sync"

	"blockchain/internal/persistence"
)

// Blockchain represents the chain and a LevelDB connection
type Blockchain struct {
	mu    sync.Mutex
	Chain []Block
	DB    *persistence.DB
}

// NewBlockchain initializes a blockchain, loading from DB if available
func NewBlockchain(db *persistence.DB) *Blockchain {
	bc := &Blockchain{
		DB: db,
	}

	// Try to load the last block from DB
	var lastBlock Block
	err := db.Get("lastBlock", &lastBlock)
	if err != nil {
		// No block in DB, start with genesis
		genesis := Genesis()
		bc.Chain = []Block{genesis}
		db.Put(genesis.Hash, genesis)
		db.Put("lastBlock", genesis)
	} else {
		// Start chain from last block
		bc.Chain = []Block{lastBlock}
	}

	return bc
}

// AddBlock mines a new block and persists it
func (bc *Blockchain) AddBlock(txs []Transaction) Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	lastBlock := bc.Chain[len(bc.Chain)-1]
	block := MineBlock(lastBlock, txs)
	bc.Chain = append(bc.Chain, block)

	// Persist the new block and update lastBlock
	bc.DB.Put(block.Hash, block)
	bc.DB.Put("lastBlock", block)

	// Optionally, persist individual transactions
	for _, tx := range txs {
		bc.DB.Put(tx.ID, tx)
	}

	return block
}

// ReplaceChain replaces the chain if the new one is longer and persists it
func (bc *Blockchain) ReplaceChain(newChain []Block) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(newChain) > len(bc.Chain) {
		bc.Chain = newChain

		// Persist all blocks
		for _, block := range newChain {
			bc.DB.Put(block.Hash, block)
		}

		// Update last block
		bc.DB.Put("lastBlock", newChain[len(newChain)-1])
	}
}
