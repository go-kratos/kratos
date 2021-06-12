package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-kratos/etcd/registry"
	pb "github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	etcd "go.etcd.io/etcd/client/v3"
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
	client, err := etcd.New(etcd.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}
	r := registry.New(client)

	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
	)
	s := &server{}
	pb.RegisterGreeterServer(grpcSrv, s)

	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			grpcSrv,
		),
		kratos.Registrar(r),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
