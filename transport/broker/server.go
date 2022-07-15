package broker

import (
	"context"

	"github.com/go-kratos/kratos/v2/broker"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

// ServerOption is an message queue server option.
type ServerOption func(*Server)

// Server is an message queue server wrapper.
type Server struct {
	broker         broker.Broker
	consumeHandler broker.ConsumeHandler
}

func NewServer(brk broker.Broker, opts ...ServerOption) transport.Server {
	srv := &Server{
		broker: brk,
	}
	for _, o := range opts {
		o(srv)
	}
	return srv
}

// Start a consumer of the message queue server
func (s *Server) Start(ctx context.Context) error {
	log.Info("Starting message queue consumer...")
	for {
		s.broker.Consume(ctx, s.consumeHandler)
	}
}

// Stop a message queue
func (s *Server) Stop(ctx context.Context) error {
	if err := s.broker.Close(); err != nil {
		return err
	}

	log.Info("Stopped message queue")
	return nil
}
