package app

import (
	"context"
	"fmt"

	signersTypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/internal/types"
	aclTypes "github.com/avareum/avareum-hubble-signer/pkg/acl/types"
	smTypes "github.com/avareum/avareum-hubble-signer/pkg/secret_manager/types"
)

type AppSigner struct {
	Signers    map[string]signersTypes.Signer
	acl        aclTypes.ACL
	sm         smTypes.SecretManager
	reqHandler types.SignerRequestedResponseHandler
}

func NewAppSigner() *AppSigner {
	a := &AppSigner{
		Signers: make(map[string]signersTypes.Signer),
	}
	return a
}

func (a *AppSigner) RegisterSecretManager(sm smTypes.SecretManager) {
	a.sm = sm
}

func (a *AppSigner) RegisterACL(acl aclTypes.ACL) {
	a.acl = acl
}

func (a *AppSigner) RegisterSignerRequestedResponseHandler(handler types.SignerRequestedResponseHandler) {
	a.reqHandler = handler
}

func (a *AppSigner) AddSigners(signers ...signersTypes.Signer) error {
	for _, s := range signers {
		err := s.Init()
		if err != nil {
			return err
		}
		a.Signers[s.ID()] = s
	}
	return nil
}

func (a *AppSigner) response(response types.SignerRequestedResponse) {
	if a.reqHandler == nil {
		return
	}

	// prevent deadlock
	select {
	case a.reqHandler <- response:
	default:
	}
}

func (a *AppSigner) Receive(ctx context.Context, mq types.MessageQueue) error {
	// register secret manager to all signer
	if a.sm == nil {
		panic("secret manager is not registered")
	}
	for _, s := range a.Signers {
		s.WithSecretManager(a.sm)
	}

	// initiate message queue connection
	receiver := mq.ReceiveChannel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case req := <-receiver:
			// check if the caller is whitelisted
			if a.acl != nil && !a.acl.CanCall(req.Caller, req.Payload, req.Signature) {
				a.response(types.SignerRequestedResponse{
					Request: req,
					Error:   fmt.Errorf("caller is not whitelisted"),
				})
				continue
			}

			if signer, isExists := a.Signers[req.SignerID()]; isExists {
				sigs, err := signer.SignAndBroadcast(ctx, req)
				if err != nil {
					a.response(types.SignerRequestedResponse{
						Request: req,
						Error:   err,
					})
				} else {
					a.response(types.SignerRequestedResponse{
						Request:    req,
						Signatures: &sigs,
					})
				}
			} else {
				a.response(types.SignerRequestedResponse{
					Request: req,
					Error:   fmt.Errorf("signer '%s' not found", req.SignerID()),
				})
			}

			// TODO: publish broadcasted signatures
		}
	}
}
