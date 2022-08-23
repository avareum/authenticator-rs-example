package types

import "fmt"

type SecretItem struct {
	prefix string
	id     string
}

func NewSecretItem(id string) SecretItem {
	return SecretItem{
		prefix: "WALLET_",
		id:     id,
	}
}

func NewSecretService(id string) SecretItem {
	return SecretItem{
		prefix: "SERVICE_",
		id:     id,
	}
}

func (w *SecretItem) ID() string {
	return fmt.Sprintf("%s%s", w.prefix, w.id)
}
