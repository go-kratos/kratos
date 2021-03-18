package queue

import "context"

// Header is a customized user header fields of the key-value pairs.
type Header map[string]string

// Message is an absctraction for all messages that
// are sent to quque or received from queue.
type Message interface {
	ID() string
	Key() string
	Header() Header
	Payload() string
	Topic() string
	Ack() error
	Nack() error
}

// Handler is a callback function that processes messages delivered
// to asynchronous subscribers.
type Handler func(Message) error

// Publisher is absctraction for sending messages
// to queue.
type Publisher interface {
	Publish(ctx context.Context, key string, payload []byte, opts ...PublishOption) error
	PublishAsync(ctx context.Context, key string, payload []byte, callback func(err error), opts ...PublishOption)
}

// Subscriber is an absctraction for receiving messages
// from queue.
type Subscriber interface {
	Subscribe(h Handler, opts ...SubscribeOption) error
	Unsubscribe() error
}
