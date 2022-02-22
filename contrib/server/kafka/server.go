package kafka

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"

	"golang.org/x/sync/errgroup"
)

var _ transport.Server = (*Server)(nil)

// Message is a message abstraction
type Message interface {
	Key() string
	Value() []byte
}

// Handler is used to handle the received message
type Handler interface {
	Topic() string
	Handle(message Message) error
}

// Consumer receives and uses the handler to process messages
type Consumer interface {
	Topic() string
	RegisterHandler(handler Handler)
	RegisterErrorHandler(handler Handler)
	Consume(ctx context.Context) error
	Close() error
}

// Server is a Kafka server wrapper
type Server struct {
	consumers     []Consumer
	handlers      map[string]Handler
	errorHandlers map[string]Handler
	logger        log.Helper
}

// ServerOption is a Kafka server option.
type ServerOption func(server *Server)

// ErrorHandlers registers a set of errorHandlers to the Server.
func ErrorHandlers(errorHandlers []Handler) ServerOption {
	return func(server *Server) {
		for _, handler := range errorHandlers {
			server.errorHandlers[handler.Topic()] = handler
		}
	}
}

// NewServer creates a Kafka server by options.
func NewServer(consumers []Consumer, handlers []Handler, opts ...ServerOption) (transport.Server, error) {
	if len(consumers) == 0 {
		return nil, fmt.Errorf("no consumers")
	}
	if len(handlers) == 0 {
		return nil, fmt.Errorf("no handlers")
	}

	server := &Server{
		consumers:     consumers,
		handlers:      make(map[string]Handler),
		errorHandlers: make(map[string]Handler),
	}

	for _, opt := range opts {
		opt(server)
	}

	for _, handler := range handlers {
		server.handlers[handler.Topic()] = handler
	}

	for _, consumer := range server.consumers {
		handler, ok := server.handlers[consumer.Topic()]
		if !ok {
			return nil, fmt.Errorf("consumer for topic %s has no handler", consumer.Topic())
		}
		consumer.RegisterHandler(handler)

		if errorHandler, ok := server.errorHandlers[consumer.Topic()]; ok {
			consumer.RegisterErrorHandler(errorHandler)
		}
	}

	return server, nil
}

// Start starts the Kafka server
func (s *Server) Start(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	for _, consumer := range s.consumers {
		consumer := consumer
		eg.Go(func() error {
			return consumer.Consume(ctx)
		})
	}

	return eg.Wait()
}

// Stop stops the Kafka server
func (s *Server) Stop(ctx context.Context) error {
	var result error
	for _, consumer := range s.consumers {
		if err := consumer.Close(); err != nil {
			s.logger.Errorf("close consumer error: %v", err)
			result = err
		}
	}

	return result
}
