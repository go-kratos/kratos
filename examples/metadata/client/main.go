package main

import (
	"context"
	"fmt"

	pb "github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/metadata/builtin"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
)

func main() {
	callHTTP()
	callGRPC()
}

func callHTTP() {
	logger := log.DefaultLogger
	conn, err := transhttp.NewClient(
		context.Background(),
		transhttp.WithMiddleware(
			recovery.Recovery(),
		),
		transhttp.WithEndpoint("127.0.0.1:8000"),
	)
	if err != nil {
		panic(err)
	}
	client := pb.NewGreeterHTTPClient(conn)
	b := &builtin.Builder{}
	ctx := metadata.NewContext(context.Background(), b.Build("kratos-extra", "2233"))
	reply, err := client.SayHello(ctx, &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		panic(err)
	}
	logger.Log(log.LevelInfo, "msg", fmt.Sprintf("[http] SayHello %s\n", reply.Message))
}

func callGRPC() {
	logger := log.DefaultLogger
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("127.0.0.1:9000"),
		transgrpc.WithMiddleware(
			middleware.Chain(
				recovery.Recovery(),
			),
		),
	)
	if err != nil {
		panic(err)
	}
	client := pb.NewGreeterClient(conn)
	b := &builtin.Builder{}
	ctx := metadata.NewContext(context.Background(), b.Build("kratos-extra", "2233"))
	reply, err := client.SayHello(ctx, &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		panic(err)
	}
	logger.Log(log.LevelInfo, "msg", fmt.Sprintf("[grpc] SayHello %+v", reply))
}
