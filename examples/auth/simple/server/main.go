package main

import (
	"context"
	"log"

	"github.com/go-kratos/kratos/examples/auth/helloworld"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

type server struct {
	helloworld.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "hello from service 2"}, nil
}

func main() {
	testKey := "serviceTestKey"
	httpSrv := http.NewServer(
		http.Address(":8001"),
		http.Middleware(
			jwt.Server(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(testKey), nil
			}),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9001"),
		grpc.Middleware(
			jwt.Server(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(testKey), nil
			}),
		),
	)
	s := &server{}
	helloworld.RegisterGreeterServer(grpcSrv, s)
	helloworld.RegisterGreeterHTTPServer(httpSrv, s)
	app := kratos.New(
		kratos.Name("helloworld2"),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
