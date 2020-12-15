package main

import (
	"context"
	"fmt"

	pb "github.com/go-kratos/kratos/v2/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println(in)
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func logger() middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			fmt.Println("start")

			return h(ctx, req)
		}
	}
}

func logger2() middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			resp, err := h(ctx, req)

			fmt.Println("end")

			return resp, err
		}
	}
}

func logger3() middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			fmt.Println("111")

			return h(ctx, req)
		}
	}
}

func main() {
	s := &server{}

	srv := grpc.NewServer(":9000", grpc.ServerMiddleware(middleware.Chain(logger(), logger2())))
	srv.Use(s, logger3())

	pb.RegisterGreeterServer(srv, s)

	srv.Start(context.Background())
}
