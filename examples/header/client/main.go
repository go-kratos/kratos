package main

import (
	"context"
	"log"
	stdhttp "net/http"

	"github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	stdgrpc "google.golang.org/grpc"
	grpcmd "google.golang.org/grpc/metadata"
)

func main() {
	callHTTP()
	callGRPC()
}

func callHTTP() {
	conn, err := http.NewClient(
		context.Background(),
		http.WithEndpoint("127.0.0.1:8000"),
	)
	if err != nil {
		panic(err)
	}
	client := helloworld.NewGreeterHTTPClient(conn)
	ctx := context.Background()
	var header stdhttp.Header
	reply, err := client.SayHello(ctx, &helloworld.HelloRequest{Name: "kratos"}, http.Header(&header))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s header: %v\n", reply.Message, header)
}

func callGRPC() {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9000"),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := helloworld.NewGreeterClient(conn)
	ctx := context.Background()
	var md grpcmd.MD
	reply, err := client.SayHello(ctx, &helloworld.HelloRequest{Name: "kratos"}, stdgrpc.Header(&md))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v header: %v\n", reply, md)
}
