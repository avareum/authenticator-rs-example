package types

import "crypto/ed25519"

type ACL interface {
	Verify(pub ed25519.PublicKey, payload []byte, payloadSignature []byte) bool
	CanCall(serviceName string, payload []byte, payloadSignature []byte) bool
}
