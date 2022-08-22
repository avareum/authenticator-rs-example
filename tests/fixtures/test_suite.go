package fixtures

import (
	signerTypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
)

type TestSuite struct {
	Solana   *SolanaTestSuite
	Ethereum *EthereumTestSuite

	MessageQueue  *TestMessageQueue
	SecretManager *TestSecretManager
	ACL           *TestACL
}

func NewTestSuite() *TestSuite {
	t := &TestSuite{
		Solana:        NewSolanaTestSuite(),
		Ethereum:      NewEthereumTestSuite(),
		MessageQueue:  NewTestMessageQueue(),
		SecretManager: NewTestSecretManager(),
		ACL:           NewTestACL(),
	}
	return t
}

func (t *TestSuite) NewTestSignerRequest() signerTypes.SignerRequest {
	return signerTypes.SignerRequest{
		Chain:     "solana",
		ChainID:   "mainnet-beta",
		Caller:    "caller-service",
		Wallet:    t.Solana.Fund.PublicKey().String(),
		Payload:   []byte{},
		Signature: []byte{},
	}
}
