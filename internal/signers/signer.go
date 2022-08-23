package signers

import (
	"fmt"

	"github.com/avareum/avareum-hubble-signer/internal/signers/types"
	smTypes "github.com/avareum/avareum-hubble-signer/pkg/secret_manager/types"
)

const WALLET_PREFIX = "WALLET_"

func ToSignerWalletID(wallet string) string {
	return fmt.Sprintf("%s%s", WALLET_PREFIX, wallet)
}

// Signer implementation checked against internal/signers/types/signer.go
var _ types.Signer = (*BaseSigner)(nil)

type BaseSigner struct {
	types.Signer
	sm smTypes.SecretManager
}

func (b *BaseSigner) Init() error {
	return nil
}

func (b *BaseSigner) WithSecretManager(sm smTypes.SecretManager) {
	b.sm = sm
}

func (b *BaseSigner) FetchSignerRawKey(id string) ([]byte, error) {
	if b.sm == nil {
		return nil, fmt.Errorf("secret manager is not set")
	}
	return b.sm.Get(ToSignerWalletID(id))
}
