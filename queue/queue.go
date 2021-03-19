package queue

import (
	"context"
)

// Event given to a subscription handler for processing.
type Event interface {
	Key() string
	Value() []byte
}

// Handler is a callback function that processes messages delivered
// to asynchronous subscribers.
type Handler func(context.Context, Event) error

// Publisher is absctraction for sending messages
// to queue.
type Publisher interface {
	Publish(ctx context.Context, key string, value []byte, opts ...PublishOption) error
	PublishAsync(ctx context.Context, key string, value []byte, callback func(err error), opts ...PublishOption) error
	Close() error
}

// Subscriber is an absctraction for receiving messages
// from queue.
type Subscriber interface {
	Subscribe(ctx context.Context, h Handler, opts ...SubscribeOption) error
	Unsubscribe(ctx context.Context) error
}
