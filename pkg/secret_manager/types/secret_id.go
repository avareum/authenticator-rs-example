package types

import "fmt"

type SecretID struct {
	prefix string
	id     string
}

func NewSecretWallet(id string) SecretID {
	return SecretID{
		prefix: "WALLET_",
		id:     id,
	}
}

func NewSecretServiceID(id string) SecretID {
	return SecretID{
		prefix: "SERVICE_",
		id:     id,
	}
}

func (w *SecretID) ID() string {
	return fmt.Sprintf("%s%s", w.prefix, w.id)
}
