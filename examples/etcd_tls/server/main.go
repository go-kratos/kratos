package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-kratos/kratos/v2/transport/http"
	"log"

	"github.com/go-kratos/etcd/registry"
	pb "github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	etcd "go.etcd.io/etcd/client/v3"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("use tls:Welcome %+v!", in.Name)}, nil
}

func main() {

	cert, err := tls.LoadX509KeyPair("../cert/server.crt", "../cert/server.key")
	if err != nil {
		panic(err)
	}
	tlsConf := &tls.Config{Certificates: []tls.Certificate{cert}}

	client, err := etcd.New(etcd.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
		),
		grpc.TLSConfig(tlsConf),
	)

	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
		),
		http.TLSConfig(tlsConf),
	)

	s := &server{}
	pb.RegisterGreeterServer(grpcSrv, s)
	pb.RegisterGreeterHTTPServer(httpSrv, s)

	r := registry.New(client)
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
