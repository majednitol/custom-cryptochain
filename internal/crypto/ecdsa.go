package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

// GenerateKeyPair creates a wallet keypair
func GenerateKeyPair() (*ecdsa.PrivateKey, string) {
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pubBytes := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	return priv, hex.EncodeToString(pubBytes)
}

// Sign signs raw data bytes with the private key (NOT hex string)
func Sign(priv *ecdsa.PrivateKey, data []byte) []byte {
	hash := HashBytes(data)

	r, s, _ := ecdsa.Sign(rand.Reader, priv, hash)
	return append(r.Bytes(), s.Bytes()...)
}

// Verify verifies an ECDSA signature given public key hex, data, and signature
func Verify(pubHex string, data []byte, sig []byte) bool {
	pubBytes, _ := hex.DecodeString(pubHex)

	x := big.Int{}
	y := big.Int{}
	x.SetBytes(pubBytes[:len(pubBytes)/2])
	y.SetBytes(pubBytes[len(pubBytes)/2:])

	pub := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     &x,
		Y:     &y,
	}

	hash := HashBytes(data)

	r := big.Int{}
	s := big.Int{}
	r.SetBytes(sig[:len(sig)/2])
	s.SetBytes(sig[len(sig)/2:])

	return ecdsa.Verify(&pub, hash, &r, &s)
}
