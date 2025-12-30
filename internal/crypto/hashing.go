package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashBytes returns raw SHA-256 bytes (used for signing & verification)
func HashBytes(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}

// Hash returns hex-encoded SHA-256 (used for blocks, logging, JSON)
func Hash(data []byte) string {
	return hex.EncodeToString(HashBytes(data))
}
