package main

import (
	"context"
	"log"
	"time"

	"github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/hashicorp/consul/api"
)

func main() {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	callHTTP(client)
	callGRPC(client)
}

func callGRPC(cli *api.Client) {
	r := registry.New(cli)
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

func callHTTP(cli *api.Client) {
	r := registry.New(cli)
	conn, err := http.NewClient(
		context.Background(),
		http.WithMiddleware(
			recovery.Recovery(),
		),
		http.WithEndpoint("discovery:///helloworld"),
		http.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 250)
	client := helloworld.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos_http"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply.Message)

}
