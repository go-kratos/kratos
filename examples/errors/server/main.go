package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-kratos/kratos/examples/errors/api"
	"github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "errors"
	// Version is the version of the compiled software.
	Version = "v1.0.0"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	fmt.Println(in.Name)
	if in.Name == "empty" {
		// Respond to errors through errors.New().
		return nil, errors.New(500, "USER_NAME_EMPTY", "user name is empty")
	}
	if in.Name == "kratos" {
		// Respond to errors with proto generated code.
		return nil, api.ErrorUserNotFound("user %s not found", "kratos")
	}
	return &helloworld.HelloReply{Message: fmt.Sprintf("Hello %+v", in.Name)}, nil
}

func main() {
	s := &server{}
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
	)
	httpSrv := http.NewServer(
		http.Address(":8000"),
	)
	helloworld.RegisterGreeterServer(grpcSrv, s)
	helloworld.RegisterGreeterHTTPServer(httpSrv, s)

	app := kratos.New(
		kratos.Name(Name),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
