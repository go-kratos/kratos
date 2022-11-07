package broker

import "context"

type Event interface {
	Key() string
	Header() map[string]string
	Value() []byte
}

type Handler func(context.Context, Event) error

type Broker interface {
	Receive(topic string, handler Handler) error
	Send(ctx context.Context, msg Event) error
	Close() error
}
