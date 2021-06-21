package main

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"log"

	"github.com/go-kratos/kratos/examples/validate/api"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
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
	v1.UnimplementedExampleServiceServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) TestValidate(ctx context.Context, in *v1.Request) (*v1.Reply, error) {
	return &v1.Reply{Message: "ok"}, nil
}

func main() {
	s := &server{}
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			middleware.Chain(
				validate.Validator(),
			),
		))
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			middleware.Chain(
				validate.Validator(),
			),
		))
	v1.RegisterExampleServiceServer(grpcSrv, s)
	v1.RegisterExampleServiceHTTPServer(httpSrv,s)
	app := kratos.New(
		kratos.Name(Name),
		kratos.Server(
			grpcSrv,
			httpSrv,
		),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
