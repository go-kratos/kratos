package main

import (
	"context"
	"log"

	"github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	grpcmd "google.golang.org/grpc/metadata"
)

func main() {
	callHTTP()
	callGRPC()
}

func callHTTP() {
	conn, err := http.NewClient(
		context.Background(),
		http.WithMiddleware(
			recovery.Recovery(),
		),
		http.WithEndpoint("127.0.0.1:8000"),
	)
	if err != nil {
		panic(err)
	}
	client := helloworld.NewGreeterHTTPClient(conn)
	md := metadata.Metadata{"kratos-extra": "2233"}
	reply, err := client.SayHello(context.Background(),
		&helloworld.HelloRequest{Name: "kratos"},
		// call options
		http.Metadata(md),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply.Message)
}

func callGRPC() {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithMiddleware(
			middleware.Chain(
				recovery.Recovery(),
			),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	client := helloworld.NewGreeterClient(conn)
	ctx := grpcmd.AppendToOutgoingContext(context.Background(), "kratos-extra", "2233")
	reply, err := client.SayHello(ctx, &helloworld.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}
