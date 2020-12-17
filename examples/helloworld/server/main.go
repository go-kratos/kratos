package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-kratos/kratos/v2"
	pb "github.com/go-kratos/kratos/v2/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/middleware"
	servergrpc "github.com/go-kratos/kratos/v2/server/grpc"
	serverhttp "github.com/go-kratos/kratos/v2/server/http"
	"github.com/go-kratos/kratos/v2/transport"
	transportgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transporthttp "github.com/go-kratos/kratos/v2/transport/http"

	"google.golang.org/grpc"
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
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			fmt.Println("before")

			return handler(ctx, req)
		}
	}
}

func logger2() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			resp, err := handler(ctx, req)

			fmt.Println("after")

			return resp, err
		}
	}
}

func logger3() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			tr, ok := transport.FromContext(ctx)
			if ok {
				fmt.Printf("transport: %+v\n", tr)
			}
			h, ok := transporthttp.FromContext(ctx)
			if ok {
				fmt.Printf("http: [%s] %s\n", h.Request.Method, h.Request.URL.Path)
			}
			g, ok := transportgrpc.FromContext(ctx)
			if ok {
				fmt.Printf("grpc: %s\n", g.FullMethod)
			}

			return handler(ctx, req)
		}
	}
}

func main() {
	s := &server{}
	app := kratos.New()

	httpTransport := transporthttp.NewServer(transporthttp.ServerMiddleware(logger(), logger2()))
	httpTransport.Use(s, logger3())

	grpcTransport := transportgrpc.NewServer(transportgrpc.ServerMiddleware(logger(), logger2()))
	grpcTransport.Use(s, logger3())

	httpServer := serverhttp.NewServer("tcp", ":8000", serverhttp.ServerHandler(httpTransport))
	grpcServer := servergrpc.NewServer("tcp", ":9000", grpc.UnaryInterceptor(grpcTransport.ServeGRPC()))

	pb.RegisterGreeterServer(grpcServer, s)
	pb.RegisterGreeterHTTPServer(httpTransport, s)

	app.Append(httpServer)
	app.Append(grpcServer)

	if err := app.Run(); err != nil {
		log.Println(err)
	}
}
