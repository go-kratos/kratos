package main

import (
    "context"
    "fmt"
    "github.com/go-kratos/etcd/registry"
    pb "github.com/go-kratos/kratos/examples/helloworld/helloworld"
    "github.com/go-kratos/kratos/v2"
    "github.com/go-kratos/kratos/v2/middleware/recovery"
    "github.com/go-kratos/kratos/v2/transport/grpc"
    "github.com/go-kratos/kratos/v2/transport/http"
    clientv3 "go.etcd.io/etcd/client/v3"
    "log"
    "testing"
)

func TestCall(t *testing.T) {
    go testServer()
    cli, err := clientv3.New(clientv3.Config{
        Endpoints: []string{"127.0.0.1:2379"},
    })
    if err != nil {
        panic(err)
    }
    r := registry.New(cli)
    callGRPC(r)
    callHTTP(r)
}

func testServer() {
    client, err := clientv3.New(clientv3.Config{
        Endpoints: []string{"127.0.0.1:2379"},
    })
    if err != nil {
        log.Fatal(err)
    }

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
    pb.RegisterGreeterServer(grpcSrv, s)
    pb.RegisterGreeterHTTPServer(httpSrv, s)

    r := registry.New(client)
    app := kratos.New(
        kratos.Name("helloworld"),
        kratos.Server(
            httpSrv,
            grpcSrv,
        ),
        kratos.Registrar(r),
    )
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
// server is used to implement helloworld.GreeterServer.
type server struct {
    pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
    return &pb.HelloReply{Message: fmt.Sprintf("Welcome %+v!", in.Name)}, nil
}