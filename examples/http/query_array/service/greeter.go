package service

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/go-kratos/kratos/examples/http/query_array/hello"
)

type GreeterService struct {
	pb.UnimplementedGreeterServer
}

func NewGreeterService() *GreeterService {
	return &GreeterService{}
}

func (s *GreeterService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	v, err := json.Marshal(req.Names)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(req.Names)
	return &pb.HelloReply{Message: string(v)}, nil
}
