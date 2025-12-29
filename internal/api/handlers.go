package api

import (
	"encoding/json"
	"net/http"

	"blockchain/internal/blockchain"
)

func RegisterHandlers(
	bc *blockchain.Blockchain,
	mempool *blockchain.Mempool,
	miner *blockchain.Miner,
) {
	http.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(bc.Chain)
	})

	http.HandleFunc("/transact", func(w http.ResponseWriter, r *http.Request) {
		var tx blockchain.Transaction
		if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mempool.Add(tx)
		w.WriteHeader(http.StatusAccepted)
	})
	http.HandleFunc("/mempool", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mempool.Transactions)
	})
	http.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request) {
	block := miner.Mine() // miner.Mine() returns nil if empty

	if block == nil {
		http.Error(w, "No transactions to mine", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(block)
})

}
