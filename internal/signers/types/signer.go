package types

import (
	"context"

	"github.com/avareum/avareum-hubble-signer/internal/types"
	smTypes "github.com/avareum/avareum-hubble-signer/pkg/secret_manager/types"
)

type Signer interface {
	ID() string
	Init() error
	WithSecretManager(sm smTypes.SecretManager)
	SignAndBroadcast(ctx context.Context, req SignerRequest) ([]string, error)
}

type SignerRequest struct {
	Chain     types.Chain
	Caller    string
	Wallet    string
	Payload   []byte
	Signature []byte
}

func (s *SignerRequest) Copy() *SignerRequest {
	return &SignerRequest{
		Chain:     s.Chain,
		Caller:    s.Caller,
		Wallet:    s.Wallet,
		Payload:   s.Payload,
		Signature: s.Payload,
	}
}

type SignerRequestedResponse struct {
	Request    SignerRequest
	Signatures []string
}

type SignerRequestedResponseHandler = chan SignerRequestedResponse
