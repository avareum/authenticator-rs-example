package signers

import (
	"os"

	"github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/pkg/secret_manager"
)

const SignerKeyVersion = "3"

type BaseSigner struct {
	types.Signer
	sm *secret_manager.SecretManager
}

func (b *BaseSigner) Init() error {
	sm, err := secret_manager.NewSecretManager(secret_manager.SecretManagerConfig{
		ProjectID:  os.Getenv("GCP_PROJECT"),
		BucketName: "signers",
	})
	b.sm = sm
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseSigner) FetchSignerRawKey() ([]byte, error) {
	return b.sm.Get(SignerKeyVersion)
}
