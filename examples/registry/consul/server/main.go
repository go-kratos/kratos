package main

import (
	"context"
	"fmt"
	"os"

	"github.com/SeeMusic/kratos/contrib/registry/consul/v2"
	"github.com/SeeMusic/kratos/examples/helloworld/helloworld"
	"github.com/SeeMusic/kratos/v2"
	"github.com/SeeMusic/kratos/v2/log"
	"github.com/SeeMusic/kratos/v2/middleware/logging"
	"github.com/SeeMusic/kratos/v2/middleware/recovery"
	"github.com/SeeMusic/kratos/v2/transport/grpc"
	"github.com/SeeMusic/kratos/v2/transport/http"
	"github.com/hashicorp/consul/api"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: fmt.Sprintf("Welcome %+v!", in.Name)}, nil
}

func main() {
	logger := log.NewStdLogger(os.Stdout)
	log := log.NewHelper(logger)

	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	)

	s := &server{}
	helloworld.RegisterGreeterServer(grpcSrv, s)
	helloworld.RegisterGreeterHTTPServer(httpSrv, s)

	r := consul.New(consulClient)
	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			grpcSrv,
			httpSrv,
		),
		kratos.Registrar(r),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
