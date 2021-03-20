package event

import (
	"context"
)

// Message is an absctraction for all messages that
// are sent to quque or received from queue.
type Message struct {
	Key    string
	Value  []byte
	Header map[string]string
}

// Event given to a subscription handler for processing.
type Event interface {
	Message() *Message
	Ack() error
	Nack() error
}

// Handler is a callback function that processes messages delivered
// to asynchronous subscribers.
type Handler func(context.Context, Event) error

// Publisher is absctraction for sending messages
// to queue.
type Publisher interface {
	Publish(ctx context.Context, msg *Message) error
	PublishAsync(ctx context.Context, msg *Message, callback func(err error)) error
	Close() error
}

// Subscriber is an absctraction for receiving messages
// from queue.
type Subscriber interface {
	Subscribe(ctx context.Context, h Handler) error
	Close() error
}
