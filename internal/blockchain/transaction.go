package blockchain

import (
	"crypto/ecdsa"
	"encoding/json"

	"blockchain/internal/crypto"
	"blockchain/internal/persistence"
)

// TxInput represents a transaction input
type TxInput struct {
	TxID      string
	Index     int
	Signature []byte
	PubKey    string
}

// TxOutput represents a transaction output
type TxOutput struct {
	Address string
	Amount  int
}

// Transaction represents a blockchain transaction
type Transaction struct {
	ID      string
	Inputs  []TxInput
	Outputs []TxOutput
}

// Sign signs the transaction with the given private key
func (tx *Transaction) Sign(privKey *ecdsa.PrivateKey) error {
	// Serialize outputs to sign
	data, err := json.Marshal(tx.Outputs)
	if err != nil {
		return err
	}

	// Sign with raw private key
	signature := crypto.Sign(privKey, data)

	// Assign signature to first input
	if len(tx.Inputs) > 0 {
		tx.Inputs[0].Signature = signature
	}

	return nil
}

// Persist saves the transaction to LevelDB (optional helper)
func (tx *Transaction) Persist(db *persistence.DB) error {
	if db == nil {
		return nil
	}
	return db.Put(tx.ID, tx)
}
