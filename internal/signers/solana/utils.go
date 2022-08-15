package solana

import (
	"encoding/base64"
	"errors"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

func Base64ToTransaction(b64 string) (*solana.Transaction, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	tx, err := solana.TransactionFromDecoder(bin.NewBinDecoder(data))
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func ProgramID(tx solana.Transaction, instructionIdx uint16) (*solana.PublicKey, error) {
	if int(instructionIdx) < len(tx.Message.Instructions) {
		return &tx.Message.AccountKeys[tx.Message.Instructions[instructionIdx].ProgramIDIndex], nil
	}
	return nil, errors.New("invalid program index")
}
