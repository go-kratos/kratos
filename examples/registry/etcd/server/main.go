package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-kratos/etcd/registry"
	pb "github.com/go-kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	clientv3 "go.etcd.io/etcd/client/v3"
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

	etcd, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}
	r := registry.New(etcd)
	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			grpcSrv,
		),
		kratos.Registrar(r),
	)

	if err := app.Run(); err != nil {
		log.Error(err)
	}
}
