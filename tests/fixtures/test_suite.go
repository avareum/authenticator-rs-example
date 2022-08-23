package fixtures

import (
	signerTypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/internal/types"
)

type TestSuite struct {
	Solana   *SolanaTestSuite
	Ethereum *EthereumTestSuite

	SecretManager *TestSecretManager
	ACL           *TestACL
}

func NewTestSuite() *TestSuite {
	t := &TestSuite{
		Solana:        NewSolanaTestSuite(),
		Ethereum:      NewEthereumTestSuite(),
		SecretManager: NewTestSecretManager(),
		ACL:           NewTestACL(),
	}
	return t
}

func (t *TestSuite) NewTestSignerRequest() signerTypes.SignerRequest {
	return signerTypes.SignerRequest{
		Chain:     types.NewChain("solana", "mainnet-beta"),
		Caller:    "caller-service",
		Wallet:    t.Solana.Fund.PublicKey().String(),
		Payload:   []byte{},
		Signature: []byte{},
	}
}
