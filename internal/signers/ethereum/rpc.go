package ethereum

import (
	"github.com/avareum/avareum-hubble-signer/internal/signers/signer"
	"github.com/avareum/avareum-hubble-signer/internal/signers/types"
)

type EthereumSignerOptions struct {
	RPC string
}

// Signer implementation checked against internal/signers/types/signer.go
var _ types.Signer = &EthereumSigner{}

type EthereumSigner struct {
	signer.BaseSigner
	opt EthereumSignerOptions
}

func NewEthereumSigner(opt EthereumSignerOptions) *EthereumSigner {
	s := &EthereumSigner{
		opt: opt,
	}
	return s
}

func (s *EthereumSigner) ID() string {
	return "ethereum.1"
}
