package main

import (
	"log"
	"net/http"

	"blockchain/internal/api"
	"blockchain/internal/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()
	mempool := blockchain.NewMempool() 
		miner := &blockchain.Miner{
		Chain:   bc,
		Mempool: mempool,
	}


	api.RegisterHandlers(bc, mempool, miner) 

	log.Println(" Node running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
