package fixtures

import (
	signerTypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/internal/types"
)

type TestMessageQueue struct {
	types.MessageQueue
	ch chan signerTypes.SignerRequest
}

func NewTestMessageQueue() *TestMessageQueue {
	return &TestMessageQueue{
		ch: make(chan signerTypes.SignerRequest),
	}
}

func (t *TestMessageQueue) Push(request signerTypes.SignerRequest) {
	t.ch <- request
}

/*
 MessageQueue implementation /types/message_queue.go
*/

func (t *TestMessageQueue) ReceiveChannel() chan signerTypes.SignerRequest {
	return t.ch
}
