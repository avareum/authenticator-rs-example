package types

type SecretManager interface {
	Create(sid SecretID, payload []byte) error
	Get(sid SecretID) ([]byte, error)
}
