package main

import (
	"context"
	"log"
	"time"

	"github.com/SeeMusic/kratos/contrib/registry/polaris/v2"

	"github.com/SeeMusic/kratos/examples/helloworld/helloworld"
	"github.com/SeeMusic/kratos/v2/registry"
	"github.com/SeeMusic/kratos/v2/transport/grpc"
	"github.com/SeeMusic/kratos/v2/transport/http"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"
)

func main() {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})
	provider, err := api.NewProviderAPIByConfig(conf)
	if err != nil {
		panic(err)
	}
	consumer, err := api.NewConsumerAPIByConfig(conf)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	defer consumer.Destroy()
	defer provider.Destroy()

	r := polaris.NewRegistry(
		provider,
		consumer,
		polaris.WithTimeout(10000),
		polaris.WithRetryCount(3),
	)
	callHTTP(r)
	callGRPC(r)
}

func callGRPC(r registry.Discovery) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworldgrpc"),
		grpc.WithDiscovery(r),
		grpc.WithTimeout(100*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}

func callHTTP(r registry.Discovery) {
	conn, err := http.NewClient(
		context.Background(),
		http.WithEndpoint("discovery:///helloworldhttp"),
		http.WithDiscovery(r),
		http.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %+v\n", reply)
}
