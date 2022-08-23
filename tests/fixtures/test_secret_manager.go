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

func (m *TestSecretManager) Create(sid types.SecretID, payload []byte) error {
	m.keystores[sid.ID()] = payload
	return nil
}

func (m *TestSecretManager) Get(sid types.SecretID) ([]byte, error) {
	rawKey, found := m.keystores[sid.ID()]
	if !found {
		return nil, fmt.Errorf("key %s not found", sid.ID())
	}
	return rawKey, nil
}
