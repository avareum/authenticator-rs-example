package solana

import (
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type SolanaTransactionDecoder struct{}

func NewSolanaTransactionDecoder() *SolanaTransactionDecoder {
	return &SolanaTransactionDecoder{}
}

func (d *SolanaTransactionDecoder) TryDecode(payload []byte) (*solana.Transaction, error) {
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
	return nil, fmt.Errorf("SolanaTransactionDecoder: unmarshal tx msg failed: %v", err)
}

func (d *SolanaTransactionDecoder) DecodeFromTransaction(payload []byte) (*solana.Transaction, error) {
	tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(payload))
	if err != nil {
		return nil, fmt.Errorf("SolanaTransactionDecoder: decode transaction failed: %v", err)
	}
	return tx, nil
}

func (d *SolanaTransactionDecoder) DecodeFromBinary(payload []byte) (*solana.Transaction, error) {
	message := solana.Message{}
	err := bin.UnmarshalBin(&message, payload)
	if err != nil {
		return nil, fmt.Errorf("SolanaTransactionDecoder: unmarshal tx msg failed: %v", err)
	}
	tx := solana.Transaction{}
	tx.Message = message
	return &tx, nil
}
