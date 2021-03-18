package queue

import (
	"context"
)

// Message is an absctraction for all messages that
// are sent to quque or received from queue.
type Message struct {
	Key   string
	Value []byte
}

// Event given to a subscription handler for processing.
type Event interface {
	Topic() string
	Message() *Message
	Ack() error
	Nack() error
}

// Handler is a callback function that processes messages delivered
// to asynchronous subscribers.
type Handler func(Event) error

// Publisher is absctraction for sending messages
// to queue.
type Publisher interface {
	Publish(ctx context.Context, msg *Message, opts ...PublishOption) error
	PublishAsync(ctx context.Context, msg *Message, callback func(err error), opts ...PublishOption) error
}

// Subscriber is an absctraction for receiving messages
// from queue.
type Subscriber interface {
	Subscribe(ctx context.Context, h Handler, opts ...SubscribeOption) error
	Unsubscribe(ctx context.Context) error
}
