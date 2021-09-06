package event

import "context"

type Event interface {
	Key() string
	Value() []byte
}

type Handler func(context.Context, Event) error

type Sender interface {
	Send(ctx context.Context, msg Event) error
	Close() error
}

type Receiver interface {
	Receive(ctx context.Context, handler Handler) error
	Close() error
}
