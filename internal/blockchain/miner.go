package blockchain

// Miner is responsible for mining blocks from the mempool
type Miner struct {
	Chain   *Blockchain
	Mempool *Mempool
}

// Mine mines transactions from the mempool, validates them, adds a new block, and updates the UTXO set
func (m *Miner) Mine(utxo *UTXOSet) *Block {
	// Get and clear pending transactions from mempool
	txs := m.Mempool.Flush()

	if len(txs) == 0 {
		return nil // Nothing to mine
	}

	validTxs := []Transaction{}

	// Validate each transaction against the UTXO set
	for _, tx := range txs {
		if err := ValidateTransaction(tx, utxo); err == nil {
			validTxs = append(validTxs, tx)
		}
	}

	if len(validTxs) == 0 {
		return nil // No valid transactions
	}

	// Add the valid transactions as a new block
	block := m.Chain.AddBlock(validTxs)

	// Update UTXO set
	for _, tx := range validTxs {
		// Remove spent UTXOs
		for _, in := range tx.Inputs {
			delete(utxo.UTXOs, in.TxID)
		}
		// Add new outputs to UTXO set
		for i, out := range tx.Outputs {
			utxo.UTXOs[tx.ID] = UTXO{
				TxID:   tx.ID,
				Index:  i,
				Output: out,
			}
		}
	}

	return &block
}
