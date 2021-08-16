package service

import (
	"context"

	pb "github.com/go-kratos/kratos/examples/http/query_array/hello"
)

type GreeterService struct {
	pb.UnimplementedGreeterServer
}

func NewGreeterService() *GreeterService {
	return &GreeterService{}
}

func (s *GreeterService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: req.Names}, nil
}
