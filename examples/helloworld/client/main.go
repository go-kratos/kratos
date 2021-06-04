package main

import (
	"context"
	"fmt"

	pb "github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
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
			logging.Client(logger),
		),
		transhttp.WithEndpoint("127.0.0.1:8000"),
	)
	if err != nil {
		panic(err)
	}
	client := pb.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		panic(err)
	}
	logger.Log(log.LevelInfo, fmt.Sprintf("[http] SayHello %s\n", reply.Message))

	// returns error
	reply, err = client.SayHello(context.Background(), &pb.HelloRequest{Name: "error"})
	if err != nil {
		logger.Log(log.LevelInfo, fmt.Sprintf("[http] SayHello error: %v\n", err))
	}
	if errors.IsBadRequest(err) {
		logger.Log(log.LevelInfo, fmt.Sprintf("[http] SayHello error is invalid argument: %v\n", err))
	}
}

func callGRPC() {
	logger := log.DefaultLogger
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("127.0.0.1:9000"),
		transgrpc.WithMiddleware(
			middleware.Chain(
				recovery.Recovery(),
				logging.Client(logger),
			),
		),
	)
	if err != nil {
		panic(err)
	}
	client := pb.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		panic(err)
	}
	logger.Log(log.LevelInfo, fmt.Sprintf("[grpc] SayHello %+v", reply))

	// returns error
	_, err = client.SayHello(context.Background(), &pb.HelloRequest{Name: "error"})
	if err != nil {
		logger.Log(log.LevelInfo, fmt.Sprintf("[grpc] SayHello error: %v", err))
	}
	if errors.IsBadRequest(err) {
		logger.Log(log.LevelInfo, fmt.Sprintf("[grpc] SayHello error is invalid argument: %v", err))
	}
}
