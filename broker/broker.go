package broker

import (
	"context"
	"io"
)

// ConsumeHandler represents a message processing method
type ConsumeHandler func(context.Context, *Message) error

type Broker interface {
	io.Closer
	Publish(context.Context, *Message) error
	Consume(context.Context, ConsumeHandler)
}
