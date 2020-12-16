package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-kratos/kratos/v2"
	pb "github.com/go-kratos/kratos/v2/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/middleware"
	transportgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transporthttp "github.com/go-kratos/kratos/v2/transport/http"

	"google.golang.org/grpc"

	_ "github.com/go-kratos/kratos/v2/encoding/json"
	_ "github.com/go-kratos/kratos/v2/encoding/proto"
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
	app := kratos.New()

	httpSrv := transporthttp.NewServer(transporthttp.ServerMiddleware(middleware.Chain(logger(), logger2())))
	httpSrv.Use(s, logger3())

	grpcSrv := transportgrpc.NewServer(transportgrpc.ServerMiddleware(middleware.Chain(logger(), logger2())))
	grpcSrv.Use(s, logger3())

	baseHTTPServer := &http.Server{Addr: ":8000", Handler: httpSrv}
	baseGRPCServer := grpc.NewServer(grpc.UnaryInterceptor(grpcSrv.ServeGRPC()))
	app.Append(kratos.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", ":9000")
			if err != nil {
				return err
			}
			pb.RegisterGreeterHTTPServer(httpSrv, s)
			return baseHTTPServer.Serve(lis)
		},
		OnStop: func(ctx context.Context) error {
			return baseHTTPServer.Shutdown(ctx)
		},
	})
	app.Append(kratos.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", ":9000")
			if err != nil {
				return err
			}
			pb.RegisterGreeterServer(baseGRPCServer, s)
			return baseGRPCServer.Serve(lis)
		},
		OnStop: func(ctx context.Context) error {
			baseGRPCServer.GracefulStop()
			return nil
		},
	})

	if err := app.Run(); err != nil {
		log.Println(err)
	}
}
