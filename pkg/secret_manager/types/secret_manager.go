package types

type SecretManager interface {
	Create(id string, payload []byte) (string, error)
	Get(id string) ([]byte, error)
}
