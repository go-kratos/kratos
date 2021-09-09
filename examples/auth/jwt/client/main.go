package main

import (
	"context"
	"log"

	"github.com/go-kratos/kratos/examples/helloworld/helloworld"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwtv4 "github.com/golang-jwt/jwt/v4"
)

type server struct {
	helloworld.UnimplementedGreeterServer

	pc helloworld.GreeterClient
}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return s.pc.SayHello(ctx, in)
}

type tokenProvider struct {
	accessSecretKey string
}

func (t tokenProvider) Key() []byte {
	return []byte(t.accessSecretKey)
}

func main() {
	testKey := "testKey"
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			jwt.Server(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(testKey), nil
			}),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			jwt.Server(func(token *jwtv4.Token) (interface{}, error) {
				return []byte(testKey), nil
			}),
		),
	)
	serviceTestKey := "serviceTestKey"
	con, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("dns:///127.0.0.1:9001"),
		grpc.WithMiddleware(
			jwt.Client(tokenProvider{
				accessSecretKey: serviceTestKey,
			}),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	s := &server{
		pc: helloworld.NewGreeterClient(con),
	}
	helloworld.RegisterGreeterServer(grpcSrv, s)
	helloworld.RegisterGreeterHTTPServer(httpSrv, s)
	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
