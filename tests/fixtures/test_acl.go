package fixtures

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"

	"github.com/avareum/avareum-hubble-signer/pkg/acl/types"
	"github.com/gagliardetto/solana-go"
)

type TestACL struct {
	types.ACL
	keypairs map[string]solana.PrivateKey
}

func NewTestACL() *TestACL {
	return &TestACL{
		keypairs: map[string]solana.PrivateKey{},
	}
}

func (a *TestACL) CreateTestServiceKey(serviceName string) error {
	servicePrivateKey, _ := solana.NewRandomPrivateKey()
	a.keypairs[serviceName] = servicePrivateKey
	return nil
}

func (a *TestACL) SignPayload(serviceName string, payload []byte) ([]byte, error) {
	priv := a.keypairs[serviceName]
	p := ed25519.PrivateKey(priv)
	signData, err := p.Sign(rand.Reader, payload, crypto.Hash(0))
	if err != nil {
		return nil, err
	}
	return signData, err
}

/*
 ACL implementaiton for GCP Secret Manager
*/

func (w *TestACL) Verify(pub ed25519.PublicKey, payload []byte, payloadSignature []byte) bool {
	return ed25519.Verify(pub, payload, payloadSignature)
}

func (w *TestACL) CanCall(serviceName string, payload []byte, payloadSignature []byte) bool {
	pub := w.keypairs[serviceName].PublicKey()
	return w.Verify(pub[:], payload, payloadSignature)
}
