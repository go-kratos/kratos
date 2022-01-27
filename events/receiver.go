package events

import "context"

type Receiver interface {
	Receive(ctx context.Context) (Message, error)
	Ack(msg Message) error
	Nack(msg Message) error
	Close() error
}

type ReceiverBuilder interface {
	Build(subReq SubRequest) (Receiver, error)
}
