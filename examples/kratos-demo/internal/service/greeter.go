package service

import (
	"fmt"

	v1 "github.com/go-kratos/kratos/v2/examples/kratos-demo/api/kratos/demo/v1"
)

// GreeterService is a Greeter service example.
type GreeterService struct {
}

// NewGreeterService new a greeter service and returns.
func NewGreeterService() *GreeterService {
	return &GreeterService{}
}

// SayHello say hello.
func (s *GreeterService) SayHello(req *v1.HelloRequest) (*v1.HelloReply, error) {
	return &v1.HelloReply{
		Message: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}
