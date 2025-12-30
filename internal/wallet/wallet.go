package wallet

import (
	"encoding/json"

	"crypto/ecdsa"

	"blockchain/internal/crypto"
	"blockchain/internal/persistence"
)

// Wallet represents a simple wallet with private key and address
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey `json:"-"`
	Address    string            `json:"address"`
}

// NewWallet generates a new wallet
func NewWallet() *Wallet {
	priv, addr := crypto.GenerateKeyPair()
	return &Wallet{
		PrivateKey: priv,
		Address:    addr,
	}
}

// Sign signs arbitrary data using the wallet's private key
func (w *Wallet) Sign(data interface{}) []byte {
	bytes, _ := json.Marshal(data)
	return crypto.Sign(w.PrivateKey, bytes)
}

// Persist saves the wallet's address to LevelDB (optional, safe)
func (w *Wallet) Persist(db *persistence.DB) error {
	if db == nil {
		return nil
	}
	return db.Put(w.Address, w)
}

// LoadWallet loads a wallet's address from LevelDB (private key must be restored separately)
func LoadWallet(db *persistence.DB, address string) (*Wallet, error) {
	if db == nil {
		return nil, nil
	}

	var w Wallet
	err := db.Get(address, &w)
	if err != nil {
		return nil, err
	}

	return &w, nil
}
