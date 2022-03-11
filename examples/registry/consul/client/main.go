package main

import (
	"context"
	"log"
	"time"

	"github.com/SeeMusic/kratos/contrib/registry/consul/v2"
	"github.com/SeeMusic/kratos/examples/helloworld/helloworld"
	"github.com/SeeMusic/kratos/v2/middleware/recovery"
	"github.com/SeeMusic/kratos/v2/transport/grpc"
	"github.com/SeeMusic/kratos/v2/transport/http"
	"github.com/hashicorp/consul/api"
)

func main() {
	consulCli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	r := consul.New(consulCli)

	// new grpc client
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	gClient := helloworld.NewGreeterClient(conn)

	// new http client
	hConn, err := http.NewClient(
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
	defer hConn.Close()
	hClient := helloworld.NewGreeterHTTPClient(hConn)

	for {
		time.Sleep(time.Second)
		callGRPC(gClient)
		callHTTP(hClient)
	}
}

func callGRPC(client helloworld.GreeterClient) {
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}

func callHTTP(client helloworld.GreeterHTTPClient) {
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply.Message)
}
