package ethereum

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type EthereumTransactionDecoder struct{}

func NewEthereumTransactionDecoder() *EthereumTransactionDecoder {
	return &EthereumTransactionDecoder{}
}

func (d *EthereumTransactionDecoder) TryDecodeFromJSON(json string) (*types.Transaction, error) {
	return d.TryDecode([]byte(json))
}

func (d *EthereumTransactionDecoder) TryDecodeFromHex(hex string) (*types.Transaction, error) {
	return d.TryDecode(common.Hex2Bytes(hex[2:]))
}

func (d *EthereumTransactionDecoder) TryDecode(payload []byte) (*types.Transaction, error) {
	// try to marshal whole tx first
	tx, err := d.DecodeFromTransaction(payload)
	if err == nil {
		return tx, nil
	}

	tx, err = d.DecodeFromJSON(payload)
	if err == nil {
		return tx, nil
	}
	return nil, fmt.Errorf("EthereumTransactionDecoder: unmarshal tx msg failed: %v", err)
}

func (d *EthereumTransactionDecoder) DecodeFromTransaction(payload []byte) (*types.Transaction, error) {
	tx := new(types.Transaction)
	err := tx.UnmarshalBinary(payload)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (d *EthereumTransactionDecoder) DecodeFromJSON(payload []byte) (*types.Transaction, error) {
	tx := new(types.Transaction)
	err := tx.UnmarshalJSON(payload)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (d *EthereumTransactionDecoder) GetMethodSignature(tx *types.Transaction) ([]byte, error) {
	return GetMethodSignature(hex.EncodeToString(tx.Data()))
}

func (d *EthereumTransactionDecoder) GetInputData(abi abi.ABI, tx *types.Transaction) ([]interface{}, error) {
	rawHex := hex.EncodeToString(tx.Data())
	sig, err := GetMethodSignature(rawHex)
	if err != nil {
		return nil, err
	}
	method, err := abi.MethodById(sig)
	if err != nil {
		return nil, err
	}
	inputData, err := GetInputData(rawHex)
	if err != nil {
		return nil, err
	}
	return method.Inputs.Unpack(inputData)
}
