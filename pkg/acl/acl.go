package acl

import (
	"crypto/ed25519"
	"fmt"

	"github.com/avareum/avareum-hubble-signer/pkg/acl/types"
	"github.com/avareum/avareum-hubble-signer/pkg/logger"
	"github.com/avareum/avareum-hubble-signer/pkg/secret_manager"
	smtypes "github.com/avareum/avareum-hubble-signer/pkg/secret_manager/types"
	"github.com/gagliardetto/solana-go"
)

type ServiceACLOptions struct {
	SkipFetchOnVerify bool
	Prefix            string
	SecretManager     smtypes.SecretManager
}

type ServiceACL struct {
	types.ACL
	opt ServiceACLOptions
	sm  smtypes.SecretManager
}

func NewServiceACL() (*ServiceACL, error) {
	return NewServiceACLWithOpt(ServiceACLOptions{
		SkipFetchOnVerify: false,
		Prefix:            "SERVICE_",
	})
}

func NewServiceACLWithOpt(opt ServiceACLOptions) (*ServiceACL, error) {
	w := &ServiceACL{
		opt: opt,
	}
	err := w.init()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *ServiceACL) init() error {
	if w.opt.SecretManager == nil {
		sm, err := secret_manager.NewGCPSecretManager()
		if err != nil {
			return err
		}
		w.sm = sm
	} else {
		w.sm = w.opt.SecretManager
	}
	return nil
}

// getServiceKey returns the public key for the given service name.
func (w *ServiceACL) getServiceKey(serviceName string) []byte {
	p, err := w.sm.Get(smtypes.NewSecretServiceID(serviceName))
	if err != nil {
		return []byte{}
	}
	return ed25519.PrivateKey(p).Public().(ed25519.PublicKey)
}

/*
 ACL implementaiton for GCP Secret Manager
*/

func (w *ServiceACL) Verify(pub ed25519.PublicKey, payload []byte, payloadSignature []byte) bool {
	return ed25519.Verify(pub, payload, payloadSignature)
}

func (w *ServiceACL) CanCall(serviceName string, payload []byte, payloadSignature []byte) bool {
	pubBytes := w.getServiceKey(serviceName)
	if len(pubBytes) == 0 {
		logger.Default.Err(fmt.Errorf("service %s not found", serviceName))
		return false
	}
	pub := solana.PublicKeyFromBytes(pubBytes)
	return w.Verify(pub[:], payload, payloadSignature)
}
