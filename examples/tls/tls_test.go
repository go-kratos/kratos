package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"testing"
	"time"

	etcdregistry "github.com/SeeMusic/kratos/contrib/registry/etcd/v2"
	pb "github.com/SeeMusic/kratos/examples/helloworld/helloworld"
	"github.com/SeeMusic/kratos/v2"
	"github.com/SeeMusic/kratos/v2/registry"
	"github.com/SeeMusic/kratos/v2/transport/grpc"
	"github.com/SeeMusic/kratos/v2/transport/http"
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

func startServer(r registry.Registrar, c *tls.Config) (app *kratos.App, err error) {
	httpSrv := http.NewServer(http.TLSConfig(c))
	grpcSrv := grpc.NewServer(grpc.TLSConfig(c))

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

func callGRPC(t *testing.T, r registry.Discovery, c *tls.Config) {
	conn, err := grpc.Dial(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithTLSConfig(c),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("[grpc] SayHello %+v\n", reply)
}

func callHTTP(t *testing.T, r registry.Discovery, c *tls.Config) {
	conn, err := http.NewClient(
		context.Background(),
		http.WithEndpoint("discovery:///helloworld"),
		http.WithTLSConfig(c),
		http.WithDiscovery(r),
		http.WithBlock(),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "kratos"})
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
	b, err := os.ReadFile("./cert/server.crt")
	if err != nil {
		t.Fatal(err)
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		t.Fatal(err)
	}
	cert, err := tls.LoadX509KeyPair("./cert/server.crt", "./cert/server.key")
	if err != nil {
		t.Fatal(err)
	}
	tlsConf := &tls.Config{
		ServerName:   "www.kratos.com",
		RootCAs:      cp,
		Certificates: []tls.Certificate{cert},
	}
	r := etcdregistry.New(client)
	srv, err := startServer(r, nil)
	if err != nil {
		t.Fatal(err)
	}
	srvTLS, err := startServer(r, tlsConf)
	if err != nil {
		t.Fatal(err)
	}
	callHTTP(t, r, tlsConf)
	callGRPC(t, r, tlsConf)
	if srv.Stop() != nil {
		t.Errorf("srv.Stop() got error: %v", err)
	}
	if srvTLS.Stop() != nil {
		t.Errorf("srv.Stop() got error: %v", err)
	}
}
