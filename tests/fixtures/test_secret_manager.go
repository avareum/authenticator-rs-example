package fixtures

import (
	"fmt"

	"github.com/avareum/avareum-hubble-signer/pkg/secret_manager/types"
)

type TestSecretManager struct {
	types.SecretManager
	keystores map[string][]byte
}

func NewTestSecretManager() *TestSecretManager {
	return &TestSecretManager{
		keystores: make(map[string][]byte),
	}
}

func (m *TestSecretManager) Create(id string, payload []byte) (string, error) {
	m.keystores[id] = payload
	return "", nil
}

func (m *TestSecretManager) Get(id string) ([]byte, error) {
	rawKey, found := m.keystores[id]
	if !found {
		return nil, fmt.Errorf("key %s not found", id)
	}
	return rawKey, nil
}
