package config

const (
	MineRate        = 2_000
	InitialDifficulty = 3
	MiningReward    = 50
)

var GenesisBlock = map[string]interface{}{
	"timestamp":  1,
	"lastHash":  "GENESIS",
	"hash":      "GENESIS_HASH",
	"nonce":     0,
	"difficulty": InitialDifficulty,
	"data":      []interface{}{},
}
