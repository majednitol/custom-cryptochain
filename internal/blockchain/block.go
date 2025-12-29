package blockchain

import (
	"encoding/json"
	"strings"
	"time"

	"blockchain/internal/crypto"
	"blockchain/pkg/config"
)

type Block struct {
	Timestamp  int64
	LastHash   string
	Hash       string
	Nonce      int
	Difficulty int
	Data       []Transaction
}

func Genesis() Block {
	raw, _ := json.Marshal(config.GenesisBlock)

	return Block{
		Timestamp:  1,
		LastHash:   "GENESIS",
		Hash:       crypto.Hash(raw),
		Nonce:      0,
		Difficulty: config.InitialDifficulty,
		Data:       []Transaction{}, // ✅ FIX
	}
}

func MineBlock(last Block, data []Transaction) Block {
	var hash string
	nonce := 0
	timestamp := time.Now().UnixMilli()
	difficulty := last.Difficulty

	for {
		blockBytes, _ := json.Marshal(struct {
			Timestamp  int64
			LastHash   string
			Data       []Transaction
			Nonce      int
			Difficulty int
		}{
			timestamp,
			last.Hash,
			data,
			nonce,
			difficulty,
		})

		hash = crypto.Hash(blockBytes)

		// ✅ FIXED PoW condition
		if strings.HasPrefix(hash, strings.Repeat("0", difficulty)) {
			break
		}

		nonce++
	}

	return Block{
		Timestamp:  timestamp,
		LastHash:   last.Hash,
		Hash:       hash,
		Nonce:      nonce,
		Difficulty: difficulty,
		Data:       data,
	}
}