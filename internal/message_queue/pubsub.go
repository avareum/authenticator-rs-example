package message_queue

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	signersTypes "github.com/avareum/avareum-hubble-signer/internal/signers/types"
	"github.com/avareum/avareum-hubble-signer/internal/types"
)

type Pubsub struct {
	types.MessageQueue
	client         *pubsub.Client
	cfg            PubsubConfig
	requestChannel chan signersTypes.SignerRequest
}

type PubsubConfig struct {
	ProjectID      string
	SubscriptionID string
	PublishID      string
}

func NewPubsub(cfg PubsubConfig) (*Pubsub, error) {
	pubsub := &Pubsub{
		cfg:            cfg,
		requestChannel: make(chan signersTypes.SignerRequest),
	}
	err := pubsub.init()
	if err != nil {
		return nil, err
	}
	return pubsub, nil
}

func (p *Pubsub) init() error {
	client, err := pubsub.NewClient(context.TODO(), p.cfg.ProjectID)
	if err != nil {
		return fmt.Errorf("Pubsub: NewClient: %v", err)
	}
	p.client = client
	return nil
}

func (p *Pubsub) Receive() error {
	sub := p.client.Subscription(p.cfg.SubscriptionID)
	err := sub.Receive(context.TODO(), func(_ context.Context, msg *pubsub.Message) {
		// TODO: parse message to signer request
		req := signersTypes.SignerRequest{}

		// trigger processor callback
		p.requestChannel <- req

		// mark as done
		msg.Ack()
	})
	return err
}

/*
 MessageQueue implementation /types/message_queue.go
*/

func (p *Pubsub) ReceiveChannel() chan signersTypes.SignerRequest {
	return p.requestChannel
}
