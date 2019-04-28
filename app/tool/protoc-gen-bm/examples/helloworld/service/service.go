package service

import (
	"context"
	"fmt"

	"go-common/app/tool/protoc-gen-bm/examples/helloworld/api"
)

// Service .
type Service struct{}

var _ v1.HelloServer = &Service{}
var _ v1.BMHelloServer = &Service{}

// SayHello .
func (s *Service) SayHello(ctx context.Context, req *v1.HelloRequest) (*v1.HelloReply, error) {
	return &v1.HelloReply{
		Message: fmt.Sprintf("hello %s", req.Name),
	}, nil
}

// Echo .
func (s *Service) Echo(ctx context.Context, req *v1.EchoRequest) (*v1.EchoReply, error) {
	return &v1.EchoReply{Content: req.Content}, nil
}
