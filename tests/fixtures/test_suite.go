package fixtures

import (
	signerTypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/internal/types"
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

func (t *TestSuite) NewSignerRequestedResponse() chan types.SignerRequestedResponse {
	return make(chan types.SignerRequestedResponse)
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
