package signers

import (
	"fmt"

	"github.com/avareum/avareum-hubble-signer/internal/signers/types"
	smtypes "github.com/avareum/avareum-hubble-signer/pkg/secret_manager/types"
)

// Signer implementation checked against internal/signers/types/signer.go
var _ types.Signer = (*BaseSigner)(nil)

type BaseSigner struct {
	types.Signer
	sm smtypes.SecretManager
}

func (b *BaseSigner) Init() error {
	return nil
}

func (b *BaseSigner) WithSecretManager(sm smtypes.SecretManager) {
	b.sm = sm
}

func (b *BaseSigner) FetchSignerRawKey(id string) ([]byte, error) {
	if b.sm == nil {
		return nil, fmt.Errorf("secret manager is not set")
	}
	return b.sm.Get(smtypes.NewSecretWallet(id))
}
