package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/go-kratos/kratos/examples/helloworld/helloworld"
	pb "github.com/go-kratos/kratos/examples/helloworld/helloworld"

	consulregistry "github.com/go-kratos/consul/registry"
	etcdregistry "github.com/go-kratos/etcd/registry"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	consul "github.com/hashicorp/consul/api"
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

func startServer(r registry.Registrar) (app *kratos.App, err error) {
	httpSrv := http.NewServer()
	grpcSrv := grpc.NewServer()

	s := &server{}
	pb.RegisterGreeterServer(grpcSrv, s)
	pb.RegisterGreeterHTTPServer(httpSrv, s)

	app = kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
		kratos.Registrar(r),
		kratos.RegistrarTimeout(5*time.Second),
	)
	go func() {
		err = app.Run()
	}()
	time.Sleep(time.Second)
	return
}

func callGRPC(t *testing.T, r registry.Discovery) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("[grpc] SayHello %+v\n", reply)
}

func callHTTP(t *testing.T, r registry.Discovery) {
	conn, err := http.NewClient(
		context.Background(),
		http.WithEndpoint("discovery:///helloworld"),
		http.WithDiscovery(r),
		http.WithBlock(),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("[http] SayHello %+v\n", reply)
}

func TestETCD(t *testing.T) {
	client, err := etcd.New(etcd.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		t.Fatal(err)
	}
	r := etcdregistry.New(client)
	srv, err := startServer(r)
	if err != nil {
		t.Fatal(err)
	}
	callHTTP(t, r)
	callGRPC(t, r)
	srv.Stop()
}

func TestConsul(t *testing.T) {
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}
	r := consulregistry.New(client)
	srv, err := startServer(r)
	if err != nil {
		t.Fatal(err)
	}
	callHTTP(t, r)
	callGRPC(t, r)
	srv.Stop()
}
