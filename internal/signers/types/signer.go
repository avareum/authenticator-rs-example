package types

import "fmt"

type Signer interface {
	ID() string
	Init() error
	SignAndBroadcast(req SignerRequest) ([]string, error)
}

type SignerRequest struct {
	Chain     string
	ChainID   string
	Payload   []byte
	Signature []byte
	Caller    string
}

func (s *SignerRequest) SignerID() string {
	return fmt.Sprintf("%s.%s", s.Chain, s.ChainID)
}

func NewMockSignerRequest() SignerRequest {
	return SignerRequest{
		Chain:     "solana",
		ChainID:   "mainnet-beta",
		Payload:   []byte("mock"),
		Signature: []byte("mock"),
		Caller:    "core",
	}
}
