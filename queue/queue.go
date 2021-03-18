package queue

import (
	"context"
	"encoding/json"
)

// Header is a customized user header fields of the key-value pairs.
type Header map[string]string

// Message is an absctraction for all messages that
// are sent to quque or received from queue.
type Message struct {
	Header map[string]string `json:"header"`
	Key    string            `json:"key"`
	Value  json.RawMessage   `json:"body"`
}

// Event given to a subscription handler for processing.
type Event interface {
	Topic() string
	Message() Message
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
	PublishAsync(ctx context.Context, msg *Message, callback func(err error), opts ...PublishOption)
}

// Subscriber is an absctraction for receiving messages
// from queue.
type Subscriber interface {
	Subscribe(h Handler, opts ...SubscribeOption) error
	Unsubscribe() error
}
