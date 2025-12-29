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

	api.RegisterHandlers(bc, mempool) 

	log.Println(" Node running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
