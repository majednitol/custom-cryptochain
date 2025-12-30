package blockchain

import (
	"sync"

	"blockchain/internal/persistence"
)

// UTXO represents an unspent transaction output
type UTXO struct {
	TxID   string
	Index  int
	Output TxOutput
}

// UTXOSet holds all unspent outputs
type UTXOSet struct {
	mu    sync.Mutex
	UTXOs map[string]UTXO
	DB    *persistence.DB
}

// NewUTXOSet creates a new UTXO set and optionally loads from DB
func NewUTXOSet(db *persistence.DB) *UTXOSet {
	u := &UTXOSet{
		UTXOs: make(map[string]UTXO),
		DB:    db,
	}

	// Load persisted UTXOs from DB if available
	if db != nil {
		var stored map[string]UTXO
		if err := db.Get("utxos", &stored); err == nil && stored != nil {
			u.UTXOs = stored
		}
	}

	return u
}

// GetBalance calculates the balance for an address
func (u *UTXOSet) GetBalance(address string) int {
	u.mu.Lock()
	defer u.mu.Unlock()

	sum := 0
	for _, utxo := range u.UTXOs {
		if utxo.Output.Address == address {
			sum += utxo.Output.Amount
		}
	}
	return sum
}

// AddUTXO adds a new UTXO and persists the set
func (u *UTXOSet) AddUTXO(utxo UTXO) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.UTXOs[utxo.TxID] = utxo
	u.persist()
}

// RemoveUTXO removes a UTXO and persists the set
func (u *UTXOSet) RemoveUTXO(txID string) {
	u.mu.Lock()
	defer u.mu.Unlock()

	delete(u.UTXOs, txID)
	u.persist()
}

// persist saves the UTXO set to LevelDB
func (u *UTXOSet) persist() {
	if u.DB != nil {
		u.DB.Put("utxos", u.UTXOs)
	}
}
