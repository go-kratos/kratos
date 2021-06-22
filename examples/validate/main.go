package main

import (
	"context"
	"log"

	"github.com/go-kratos/kratos/examples/validate/api"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/validate"
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

type server struct {
	v1.UnimplementedExampleServiceServer
}

func (s *server) TestValidate(ctx context.Context, in *v1.Request) (*v1.Reply, error) {
	return &v1.Reply{Message: "ok"}, nil
}

func main() {
	s := &server{}
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			validate.Validator(),
		))
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			validate.Validator(),
		))
	v1.RegisterExampleServiceServer(grpcSrv, s)
	v1.RegisterExampleServiceHTTPServer(httpSrv, s)
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
