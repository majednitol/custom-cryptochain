package blockchain

import "sync"

type Mempool struct {
	mu           sync.Mutex
	Transactions []Transaction
}

func NewMempool() *Mempool {
	return &Mempool{
		Transactions: []Transaction{},
	}
}

func (m *Mempool) Add(tx Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Transactions = append(m.Transactions, tx)
}

func (m *Mempool) Flush() []Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()
	txs := m.Transactions
	if txs == nil {
		txs = []Transaction{}
	}

	m.Transactions = []Transaction{}
	return txs
}
