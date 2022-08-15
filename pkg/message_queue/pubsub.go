package message_queue

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/avareum/avareum-hubble-signer/pkg/message_queue/types"
)

type Pubsub struct {
	types.IMessageQueue
	client          *pubsub.Client
	responseChannel chan types.MessageQueueRequest
	cfg             PubsubConfig
}

type PubsubConfig struct {
	ProjectID      string
	SubscriptionID string
	PublishID      string
}

func NewPubsub(cfg PubsubConfig) *Pubsub {
	pubsub := &Pubsub{
		cfg: cfg,
	}
	pubsub.init()
	return pubsub
}

func (p *Pubsub) init() error {
	client, err := pubsub.NewClient(context.TODO(), p.cfg.ProjectID)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %v", err)
	}
	p.client = client
	return nil
}

func (p *Pubsub) Receive(callback func(message types.MessageQueueRequest) error) error {
	sub := p.client.Subscription(p.cfg.SubscriptionID)
	return sub.Receive(context.TODO(), func(_ context.Context, msg *pubsub.Message) {
		// TODO: 1. parse pubsub message to MessageQueueRequest

		// TODO: 2. trigger processor callback
		callback(types.MessageQueueRequest{})

		// 3. mark as done
		msg.Ack()
	})
}
