package types

import "github.com/avareum/avareum-hubble-signer/internal/signers/types"

type MessageQueueError struct {
	Request types.SignerRequest
	Error   error
}

type MessageQueue interface {
	ReceiveChannel() chan types.SignerRequest
}
