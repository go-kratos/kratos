package main

import (
	"context"
	"log"

	"github.com/SeeMusic/kratos/examples/helloworld/helloworld"
	"github.com/SeeMusic/kratos/v2/metadata"
	mmd "github.com/SeeMusic/kratos/v2/middleware/metadata"
	"github.com/SeeMusic/kratos/v2/transport/grpc"
	"github.com/SeeMusic/kratos/v2/transport/http"
)

func main() {
	callHTTP()
	callGRPC()
}

func callHTTP() {
	conn, err := http.NewClient(
		context.Background(),
		http.WithMiddleware(
			mmd.Client(),
		),
		http.WithEndpoint("127.0.0.1:8000"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterHTTPClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToClientContext(ctx, "x-md-global-extra", "2233")
	reply, err := client.SayHello(ctx, &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply)
}

func callGRPC() {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithMiddleware(
			mmd.Client(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToClientContext(ctx, "x-md-global-extra", "2233")
	reply, err := client.SayHello(ctx, &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v \n", reply)
}
