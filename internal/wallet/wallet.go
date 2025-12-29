package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  string
}

func NewWallet() *Wallet {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub := append(key.PublicKey.X.Bytes(), key.PublicKey.Y.Bytes()...)
	return &Wallet{
		PrivateKey: key,
		PublicKey:  hex.EncodeToString(pub),
	}
}
