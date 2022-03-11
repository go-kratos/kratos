package main

import (
	"context"
	"log"

	stdhttp "net/http"

	"github.com/SeeMusic/kratos/examples/errors/api"
	pb "github.com/SeeMusic/kratos/examples/helloworld/helloworld"
	"github.com/SeeMusic/kratos/v2/errors"
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
		http.WithEndpoint("127.0.0.1:8000"),
	)
	if err != nil {
		panic(err)
	}
	client := pb.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "empty"})
	if err != nil {
		if errors.Code(err) == stdhttp.StatusInternalServerError {
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
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9000"),
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
