package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"blockchain/internal/blockchain"
	"blockchain/internal/wallet"
)

// RegisterHandlers sets up all blockchain HTTP endpoints
func RegisterHandlers(
	bc *blockchain.Blockchain,
	mempool *blockchain.Mempool,
	miner *blockchain.Miner,
	utxo *blockchain.UTXOSet,
) {
	// -----------------------------
	// Blocks
	// -----------------------------
	http.HandleFunc("/blocks", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(bc.Chain)
	})

	http.HandleFunc("/block/", func(w http.ResponseWriter, r *http.Request) {
		hash := strings.TrimPrefix(r.URL.Path, "/block/")
		for _, block := range bc.Chain {
			if block.Hash == hash {
				json.NewEncoder(w).Encode(block)
				return
			}
		}
		http.Error(w, "Block not found", http.StatusNotFound)
	})

	// -----------------------------
	// Transactions
	// -----------------------------
	http.HandleFunc("/transact", func(w http.ResponseWriter, r *http.Request) {
		var tx blockchain.Transaction
		if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mempool.Add(tx)
		w.WriteHeader(http.StatusAccepted)
	})

	http.HandleFunc("/transaction/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/transaction/")
		// Search mempool first
		for _, tx := range mempool.Transactions {
			if tx.ID == id {
				json.NewEncoder(w).Encode(tx)
				return
			}
		}
		// Search blockchain
		for _, block := range bc.Chain {
			for _, tx := range block.Data {
				if tx.ID == id {
					json.NewEncoder(w).Encode(tx)
					return
				}
			}
		}
		http.Error(w, "Transaction not found", http.StatusNotFound)
	})

	// -----------------------------
	// Mempool
	// -----------------------------
	http.HandleFunc("/mempool", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(mempool.Transactions)
	})

	// -----------------------------
	// Mining
	// -----------------------------
	http.HandleFunc("/mine", func(w http.ResponseWriter, r *http.Request) {
		block := miner.Mine(utxo)
		if block == nil {
			http.Error(w, "No valid transactions to mine", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(block)
	})

	// -----------------------------
	// Wallets
	// -----------------------------
	http.HandleFunc("/wallet/new", func(w http.ResponseWriter, r *http.Request) {
		newWallet := wallet.NewWallet()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"address": newWallet.Address,
		})
	})

	http.HandleFunc("/wallet/send", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			From    string
			To      string
			Amount  int
			PrivKey string // optional hex string, for signing
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO: Implement wallet transaction creation & signing
		http.Error(w, "Wallet send endpoint not fully implemented", http.StatusNotImplemented)
	})

	// -----------------------------
	// Balances & UTXOs
	// -----------------------------
	http.HandleFunc("/balance/", func(w http.ResponseWriter, r *http.Request) {
		address := strings.TrimPrefix(r.URL.Path, "/balance/")
		if address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}
		balance := utxo.GetBalance(address)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"address": address,
			"balance": balance,
		})
	})

	http.HandleFunc("/utxos/", func(w http.ResponseWriter, r *http.Request) {
		address := strings.TrimPrefix(r.URL.Path, "/utxos/")
		if address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}

		var userUTXOs []blockchain.UTXO
		for _, u := range utxo.UTXOs {
			if u.Output.Address == address {
				userUTXOs = append(userUTXOs, u)
			}
		}

		json.NewEncoder(w).Encode(userUTXOs)
	})

	// -----------------------------
	// Node info (optional)
	// -----------------------------
	http.HandleFunc("/chain/length", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]int{
			"length": len(bc.Chain),
		})
	})

	// -----------------------------
	// Networking / Peers (if applicable)
	// -----------------------------
	// http.HandleFunc("/peers", ...)
	// http.HandleFunc("/addpeer", ...)
}
