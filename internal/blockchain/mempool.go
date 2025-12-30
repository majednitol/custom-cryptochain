package blockchain

import (
	"blockchain/internal/persistence"
	"sync"
)

// Mempool holds unconfirmed transactions
type Mempool struct {
	mu           sync.Mutex
	Transactions []Transaction
	DB           *persistence.DB
}

// NewMempool initializes a mempool with optional LevelDB persistence
func NewMempool(db *persistence.DB) *Mempool {
	m := &Mempool{
		Transactions: []Transaction{},
		DB:           db,
	}

	// Load existing pending transactions from DB if available
	if db != nil {
		var storedTxs []Transaction
		err := db.Get("mempool", &storedTxs)
		if err == nil && storedTxs != nil {
			m.Transactions = storedTxs
		}
	}

	return m
}

// Add adds a new transaction to the mempool and persists it
func (m *Mempool) Add(tx Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Transactions = append(m.Transactions, tx)

	// Persist the mempool
	if m.DB != nil {
		m.DB.Put("mempool", m.Transactions)
		// Also persist individual transaction by ID
		m.DB.Put(tx.ID, tx)
	}
}

// Flush clears the mempool and returns all transactions
func (m *Mempool) Flush() []Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()

	txs := m.Transactions
	if txs == nil {
		txs = []Transaction{}
	}

	// Clear the mempool in memory
	m.Transactions = []Transaction{}

	// Clear the persisted mempool
	if m.DB != nil {
		m.DB.Put("mempool", m.Transactions)
	}

	return txs
}
