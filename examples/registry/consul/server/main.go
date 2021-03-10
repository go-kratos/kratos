package main

import (
	"context"
	"fmt"
	"os"

	consul "github.com/go-kratos/consul/registry"
	pb "github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Welcome %+v!", in.Name)}, nil
}

func main() {
	logger := log.NewStdLogger(os.Stdout)

	log := log.NewHelper("main", logger)

	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
	)

	s := &server{}
	pb.RegisterGreeterServer(grpcSrv, s)

	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	r, err := consul.New(cli)
	if err != nil {
		panic(err)
	}

	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			grpcSrv,
		),
		kratos.Registry(r),
	)

	if err := app.Run(); err != nil {
		log.Error(err)
	}
}
