package ethereum

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

type EthereumTransactionDecoder struct{}

func NewEthereumTransactionDecoder() *EthereumTransactionDecoder {
	return &EthereumTransactionDecoder{}
}

func (d *EthereumTransactionDecoder) TryDecode(payload []byte) (*types.Transaction, error) {
	// try to marshal whole tx first
	tx, err := d.DecodeFromTransaction(payload)
	if err == nil {
		return tx, nil
	}

	// otherwise, try to unmarshal only tx message
	tx, err = d.DecodeFromBinary(payload)
	if err == nil {
		return tx, nil
	}
	return nil, fmt.Errorf("EthereumTransactionDecoder: unmarshal tx msg failed: %v", err)
}

func (d *EthereumTransactionDecoder) DecodeFromTransaction(payload []byte) (*types.Transaction, error) {
	tx := new(types.Transaction)
	err := tx.UnmarshalBinary(payload)
	if err != nil {
		return nil, fmt.Errorf("EthereumTransactionDecoder: unmarshal tx msg failed: %v", err)
	}
	return tx, nil
}

func (d *EthereumTransactionDecoder) DecodeFromBinary(payload []byte) (*types.Transaction, error) {
	return nil, fmt.Errorf("EthereumTransactionDecoder: DecodeFromBinary is not implemented")
}
