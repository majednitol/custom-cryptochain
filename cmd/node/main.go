package main

import (
	"log"
	"net/http"

	"blockchain/internal/api"
	"blockchain/internal/blockchain"
	"blockchain/internal/persistence"
)

func main() {
	// Open LevelDB
	db, err := persistence.Open("data/leveldb")
	if err != nil {
		log.Fatal("Failed to open LevelDB:", err)
	}

	// Initialize blockchain with DB
	bc := blockchain.NewBlockchain(db)

	// Initialize mempool and UTXO set with DB
	mempool := blockchain.NewMempool(db)
	utxo := blockchain.NewUTXOSet(db)

	// Miner
	miner := &blockchain.Miner{
		Chain:   bc,
		Mempool: mempool,
	}

	// API
	api.RegisterHandlers(bc, mempool, miner, utxo)

	log.Println("🚀 Node running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
