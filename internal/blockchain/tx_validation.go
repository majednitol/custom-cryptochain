package blockchain

import (
	"encoding/json"
	"errors"

	"blockchain/internal/crypto"
	"blockchain/internal/persistence"
)

// ValidateTransaction checks that inputs exist, signatures are valid, and sums match
func ValidateTransaction(tx Transaction, utxoSet *UTXOSet) error {
	inputSum := 0
	outputSum := 0

	for _, in := range tx.Inputs {
		utxo, ok := utxoSet.UTXOs[in.TxID]
		if !ok {
			return errors.New("UTXO not found")
		}

		// Serialize outputs to verify signature
		data, _ := json.Marshal(tx.Outputs)
		if !crypto.Verify(in.PubKey, data, in.Signature) {
			return errors.New("invalid signature")
		}

		inputSum += utxo.Output.Amount
	}

	for _, out := range tx.Outputs {
		outputSum += out.Amount
	}

	if inputSum != outputSum {
		return errors.New("input/output mismatch")
	}

	return nil
}

// Optional helper to persist validated transaction to LevelDB
func ValidateAndPersist(tx Transaction, utxoSet *UTXOSet, db *persistence.DB) error {
	if err := ValidateTransaction(tx, utxoSet); err != nil {
		return err
	}

	if db != nil {
		return db.Put(tx.ID, tx)
	}

	return nil
}
