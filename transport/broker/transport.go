package broker

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport"
)

type broker struct {
	broker   Broker
	receives map[string]Handler
}

type HandlerFunc func(context.Context, Event) error

func (b *broker) Start(_ context.Context) error {
	for topic, handler := range b.receives {
		if err := b.broker.Receive(topic, func(ctx context.Context, event Event) error {
			return handler(ctx, event)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (b *broker) Stop(ctx context.Context) error {
	return b.broker.Close()
}

func NewBroker(b Broker) transport.Server {
	return &broker{}
}

// Receive registers a handler for the given topic.
func (b *broker) Receive(topic string, handler Handler) error {
	return b.broker.Receive(topic, handler)
}
