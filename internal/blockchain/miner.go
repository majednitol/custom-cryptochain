package blockchain


type Miner struct {
	Chain   *Blockchain
	Mempool *Mempool
}

func (m *Miner) Mine() Block {
	txs := m.Mempool.Flush()
	return m.Chain.AddBlock(txs)
}
