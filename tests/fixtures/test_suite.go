package fixtures

import (
	signerTypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
)

type TestSuite struct {
	Solana        *SolanaTestSuite
	MessageQueue  *TestMessageQueue
	SecretManager *TestSecretManager
	ACL           *TestACL
}

func NewTestSuite() *TestSuite {
	t := &TestSuite{
		Solana:        NewSolanaTestSuite(),
		MessageQueue:  NewTestMessageQueue(),
		SecretManager: NewTestSecretManager(),
		ACL:           NewTestACL(),
	}
	return t
}

func (t *TestSuite) Faucet() *TestSuite {
	t.Solana.Airdrop()
	return t
}

func (t *TestSuite) NewTestSignerRequest() signerTypes.SignerRequest {
	return signerTypes.SignerRequest{
		Chain:     "solana",
		ChainID:   "mainnet-beta",
		Caller:    "caller-service",
		Fund:      t.Solana.Fund.PublicKey().String(),
		Payload:   []byte{},
		Signature: []byte{},
	}
}
