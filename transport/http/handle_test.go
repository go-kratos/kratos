package http

import (
	"context"
	"testing"
)

type HelloRequest struct {
	Name string `json:"name"`
}
type HelloReply struct {
	Message string `json:"message"`
}
type GreeterService struct {
}

func (s *GreeterService) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "hello " + req.Name}, nil
}

func TestHandler(t *testing.T) {
	s := &GreeterService{}
	_ = NewHandler(s.SayHello)
}
