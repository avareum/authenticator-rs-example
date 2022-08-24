package app

import (
	"context"
	"fmt"

	signerTypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
	aclTypes "github.com/avareum/avareum-hubble-signer/pkg/acl/types"
	smTypes "github.com/avareum/avareum-hubble-signer/pkg/secret_manager/types"
)

type AppSigner struct {
	Signers map[string]signerTypes.Signer
	acl     aclTypes.ACL
	sm      smTypes.SecretManager
}

func NewAppSigner() *AppSigner {
	a := &AppSigner{
		Signers: make(map[string]signerTypes.Signer),
	}
	return a
}

func (a *AppSigner) WithSecretManager(sm smTypes.SecretManager) *AppSigner {
	a.sm = sm
	return a
}

func (a *AppSigner) WithACL(acl aclTypes.ACL) *AppSigner {
	a.acl = acl
	return a
}

func (a *AppSigner) AddSigners(signers ...signerTypes.Signer) error {
	for _, s := range signers {
		err := s.Init()
		if err != nil {
			return err
		}
		a.Signers[s.Chain().ID()] = s
	}
	return nil
}

func (a *AppSigner) TrySign(ctx context.Context, req signerTypes.SignerRequest) (*signerTypes.SignerRequestedResponse, error) {
	// register secret manager to all signer
	if a.sm == nil {
		return nil, fmt.Errorf("secret manager is not registered")
	}
	for _, s := range a.Signers {
		s.WithSecretManager(a.sm)
	}

	// check if the caller is whitelisted
	if a.acl != nil && !a.acl.CanCall(req.Caller, req.Payload, req.Signature) {
		return nil, fmt.Errorf("invalid caller signature")
	}

	if signer, isExists := a.Signers[req.Chain.ID()]; isExists {
		sigs, err := signer.SignAndBroadcast(ctx, req)
		if err != nil {
			return nil, err
		} else {
			return &signerTypes.SignerRequestedResponse{
				Request:    req,
				Signatures: sigs,
			}, nil
		}
	} else {
		return nil, fmt.Errorf("signer %s not found", req.Chain.ID())
	}
}
