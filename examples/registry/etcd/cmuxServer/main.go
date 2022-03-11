package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/SeeMusic/kratos/contrib/registry/etcd/v2"
	pb "github.com/SeeMusic/kratos/examples/helloworld/helloworld"
	"github.com/SeeMusic/kratos/v2"
	"github.com/SeeMusic/kratos/v2/middleware/recovery"
	"github.com/SeeMusic/kratos/v2/transport/grpc"
	"github.com/SeeMusic/kratos/v2/transport/http"
	"github.com/soheilhy/cmux"
	etcdclient "go.etcd.io/etcd/client/v3"
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
	l, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Panic(err)
	}
	m := cmux.New(l)

	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}

	grpcSrv := grpc.NewServer(
		grpc.Listener(m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))),
		grpc.Middleware(
			recovery.Recovery(),
		),
	)
	httpSrv := http.NewServer(
		http.Listener(m.Match(cmux.Any())),
		http.Middleware(
			recovery.Recovery(),
		),
	)

	s := &server{}
	pb.RegisterGreeterServer(grpcSrv, s)
	pb.RegisterGreeterHTTPServer(httpSrv, s)

	r := etcd.New(client)
	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
		kratos.Registrar(r),
	)

	go func() {
		if err := m.Serve(); !strings.Contains(err.Error(), "use of closed network connection") {
			panic(err)
		}
	}()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
