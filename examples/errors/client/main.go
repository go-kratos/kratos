package main

import (
	"context"
	"github.com/go-kratos/kratos/examples/errors/api"
	"github.com/go-kratos/kratos/v2/errors"
	"log"

	pb "github.com/go-kratos/kratos/examples/helloworld/helloworld"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
)

func main() {
	callHTTP()
	callGRPC()
}

func callHTTP() {
	conn, err := transhttp.NewClient(
		context.Background(),
		transhttp.WithEndpoint("127.0.0.1:8000"),
	)
	if err != nil {
		panic(err)
	}
	client := pb.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "empty"})
	if err != nil {
		if errors.Code(err) == 500 {
			log.Println(err)
		}
		if api.IsUserNotFound(err) {
			log.Println("[http] USER_NOT_FOUND_ERROR", err)
		}
	} else {
		log.Printf("[http] SayHello %s\n", reply.Message)
	}
}

func callGRPC() {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("127.0.0.1:9000"),
	)
	if err != nil {
		panic(err)
	}
	client := pb.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		e := errors.FromError(err)
		if e.Reason == "USER_NAME_EMPTY" && e.Code == 500 {
			log.Println("[grpc] USER_NAME_EMPTY", err)
		}
		if api.IsUserNotFound(err) {
			log.Println("[grpc] USER_NOT_FOUND_ERROR", err)
		}
	} else {
		log.Printf("[grpc] SayHello %+v\n", reply)
	}
}
