package utils

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
)

// Verify checks if a message is signed by a given Public Key
func Verify(digest []byte, sig []byte, pk *rsa.PublicKey) error {
	h := sha256.New()
	h.Write(digest)
	d := h.Sum(nil)
	return rsa.VerifyPKCS1v15(pk, crypto.SHA256, d, sig)
}
