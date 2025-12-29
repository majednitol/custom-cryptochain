package blockchain

import "sync"

type Blockchain struct {
	mu    sync.Mutex
	Chain []Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Chain: []Block{Genesis()},
	}
}

// ✅ FIX: interface{} → []Transaction
func (bc *Blockchain) AddBlock(txs []Transaction) Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	block := MineBlock(bc.Chain[len(bc.Chain)-1], txs)
	bc.Chain = append(bc.Chain, block)
	return block
}

func (bc *Blockchain) ReplaceChain(newChain []Block) {
	if len(newChain) > len(bc.Chain) {
		bc.Chain = newChain
	}
}
