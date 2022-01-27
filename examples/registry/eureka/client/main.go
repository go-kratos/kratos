package main

import (
	"context"
	"log"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/eureka/v2"

	"github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	srcgrpc "google.golang.org/grpc"
)

func main() {

	r, err := eureka.New([]string{"http://127.0.0.1:18761"}, eureka.WithRefresh("1s"))

	if err != nil {
		log.Fatal(err)
	}

	connHTTP, err := http.NewClient(
		context.Background(),
		http.WithEndpoint("discovery:///helloworld"),
		http.WithDiscovery(r),
		http.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer connHTTP.Close()

	connGRPC, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer connGRPC.Close()

	for {
		callHTTP(r, connHTTP)
		callGRPC(r, connGRPC)
		time.Sleep(time.Second)
	}
}

func callGRPC(r *eureka.Registry, conn *srcgrpc.ClientConn) {
	client := helloworld.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "rocky"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}

func callHTTP(r *eureka.Registry, conn *http.Client) {
	client := helloworld.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "rocky"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %+v\n", reply)
}
