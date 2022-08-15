package secret_manager

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type SecretManager struct {
	client *secretmanager.Client
	cfg    SecretManagerConfig
}

type SecretManagerConfig struct {
	ProjectID  string
	BucketName string
}

func NewSecretManager(cfg SecretManagerConfig) (*SecretManager, error) {
	sm := &SecretManager{
		cfg: cfg,
	}
	err := sm.init()
	if err != nil {
		return nil, err
	}
	return sm, nil
}

func (s *SecretManager) init() error {
	client, err := secretmanager.NewClient(context.TODO())
	if err != nil {
		return fmt.Errorf("secretmanager: NewClient: %v", err)
	}
	s.client = client
	return nil
}

func (s *SecretManager) Create(id string, payload []byte) (string, error) {
	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", s.cfg.ProjectID),
		SecretId: id,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	secret, err := s.client.CreateSecret(context.TODO(), createSecretReq)
	if err != nil {
		return "", fmt.Errorf("failed to create secret: %v", err)
	}

	// Build the request.
	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: secret.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	// Call the API.
	version, err := s.client.AddSecretVersion(context.TODO(), addSecretVersionReq)
	if err != nil {
		return "", fmt.Errorf("failed to add secret version: %v", err)
	}

	return version.Name, nil
}

func (s *SecretManager) Get(keyVersion string) ([]byte, error) {
	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", s.cfg.ProjectID, s.cfg.BucketName, keyVersion),
	}

	// Call the API.
	result, err := s.client.AccessSecretVersion(context.TODO(), accessRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %v", err)
	}

	return result.Payload.Data, nil
}
