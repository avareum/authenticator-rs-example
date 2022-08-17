package utils

import (
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

// MarshalBinarySolanaTransaction is a helper struct to marshal a transaction without signature verification checked (NumRequiredSignatures)
func MarshalBinarySolanaTransaction(tx *solana.Transaction) ([]byte, error) {
	messageContent, err := tx.Message.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to encode tx.Message to binary: %w", err)
	}

	var signatureCount []byte
	bin.EncodeCompactU16Length(&signatureCount, len(tx.Signatures))
	output := make([]byte, 0, len(signatureCount)+len(signatureCount)*64+len(messageContent))
	output = append(output, signatureCount...)
	for _, sig := range tx.Signatures {
		output = append(output, sig[:]...)
	}
	output = append(output, messageContent...)
	return output, nil
}
