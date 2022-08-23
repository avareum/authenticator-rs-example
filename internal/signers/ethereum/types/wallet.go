package types

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func NewEthereumKey() (*ecdsa.PrivateKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func MustNewEthereumKey() *ecdsa.PrivateKey {
	privateKey, err := NewEthereumKey()
	if err != nil {
		log.Fatal(err)
	}
	return privateKey
}
