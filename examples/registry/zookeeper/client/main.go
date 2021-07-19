package main

import (
	"context"
	"log"

	"github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pokitpeng/zookeeper/registry"
)

func main() {
	r, err := registry.New([]string{"127.0.0.1:2181"})
	if err != nil {
		panic(err)
	}
	callHTTP(r)
	callGRPC(r)
}

func callGRPC(r *registry.Registry) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := helloworld.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos_grpc"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}

func callHTTP(r *registry.Registry) {
	conn, err := http.NewClient(
		context.Background(),
		http.WithMiddleware(
			recovery.Recovery(),
		),
		http.WithEndpoint("discovery:///helloworld"),
		http.WithDiscovery(r),
		http.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := helloworld.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos_http"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply)
}
