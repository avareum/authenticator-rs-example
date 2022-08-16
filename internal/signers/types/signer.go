package types

import (
	"context"
	"fmt"

	smTypes "github.com/avareum/avareum-hubble-signer/pkg/secret_manager/types"
)

type Signer interface {
	ID() string
	Init() error
	WithSecretManager(sm smTypes.SecretManager)
	SignAndBroadcast(ctx context.Context, req SignerRequest) ([]string, error)
}

type SignerRequest struct {
	Chain     string
	ChainID   string
	Caller    string
	Fund      string
	Payload   []byte
	Signature []byte
}

func (s *SignerRequest) SignerID() string {
	return fmt.Sprintf("%s.%s", s.Chain, s.ChainID)
}
