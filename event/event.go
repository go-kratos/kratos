package event

import (
	"context"
)

// Event is an absctraction for all messages that
// are sent to quque or received from queue.
type Event struct {
	// Key sets the key of the message for routing policy
	Key string
	// Payload for the message
	Payload []byte
	// Properties attach application defined properties on the message
	Properties map[string]string
}

// Handler is a callback function that processes messages delivered
// to asynchronous subscribers.
type Handler func(context.Context, Event) error

// Sender is absctraction for sending messages
// to queue.
type Sender interface {
	Send(ctx context.Context, event Event) error
	SendAsync(ctx context.Context, event Event, callback func(err error)) error
	Close() error
}

// Receiver is an absctraction for receiving messages
// from queue.
type Receiver interface {
	Receive(ctx context.Context, h Handler) error
	Close() error
}
