package main

import (
	"context"
	"fmt"
	"log"

	"github.com/SeeMusic/kratos/contrib/registry/polaris/v2"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"

	"github.com/SeeMusic/kratos/examples/helloworld/helloworld"
	"github.com/SeeMusic/kratos/v2"
	"github.com/SeeMusic/kratos/v2/middleware/recovery"
	"github.com/SeeMusic/kratos/v2/transport/grpc"
	"github.com/SeeMusic/kratos/v2/transport/http"
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
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	provider, err := api.NewProviderAPIByConfig(conf)
	if err != nil {
		panic(err)
	}

	consumer, err := api.NewConsumerAPIByConfig(conf)
	if err != nil {
		panic(err)
	}
	defer consumer.Destroy()
	defer provider.Destroy()

	registry := polaris.NewRegistry(
		provider,
		consumer,
		polaris.WithTTL(5),
	)

	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
		),
	)

	s := &server{}

	helloworld.RegisterGreeterServer(grpcSrv, s)
	helloworld.RegisterGreeterHTTPServer(httpSrv, s)

	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			grpcSrv,
			httpSrv,
		),
		kratos.Registrar(registry),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
