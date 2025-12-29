package api

import (
	"encoding/json"
	"net/http"

	"blockchain/internal/blockchain"
)

func RegisterHandlers(
	bc *blockchain.Blockchain,
	mempool *blockchain.Mempool,
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
		block := bc.AddBlock(mempool.Flush())
		json.NewEncoder(w).Encode(block)
	})
}
