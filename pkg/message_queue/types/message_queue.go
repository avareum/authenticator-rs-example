package types

type IMessageQueue interface {
	Receive(func(message MessageQueueRequest) error) error
}

type MessageQueueRequest struct {
	// TODO:
}
