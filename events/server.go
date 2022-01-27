package events

import (
	"context"
	"log"
)

type ServerOption func(*Server)

func WithReceiverBuilder(builder ReceiverBuilder) ServerOption {
	return func(s *Server) {
		s.receiverBuilder = builder
	}
}

type Server struct {
	receiverBuilder ReceiverBuilder
	receiverMap     map[Receiver]Handler
}

func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		receiverMap: make(map[Receiver]Handler),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Start(ctx context.Context) error {
	for receiver, handler := range s.receiverMap {
		go func(receiver Receiver, handler Handler) {
			for {
				msg, err := receiver.Receive(ctx)
				if err != nil {
					log.Printf("receiver error: %v", err)
					return
				}
				err = handler.Handle(context.Background(), msg)
				if err != nil {
					log.Printf("handler error: %v", err)
				}
			}
		}(receiver, handler)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	for receiver := range s.receiverMap {
		err := receiver.Close()
		if err != nil {
			log.Printf("Error closing receiver: %v", err)
		}
	}
	return nil
}

func (s *Server) Subscribe(subReq SubRequest, handler Handler) error {
	receiver, err := s.receiverBuilder.Build(subReq)
	if err != nil {
		return err
	}
	s.receiverMap[receiver] = handler
	return nil
}
