package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/grpc/health"
)

func main() {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("127.0.0.1:9000"),
		transgrpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := health.NewHealthClient(conn)
	stream, err := client.Watch(context.Background(), &health.HealthCheckRequest{Service: "helloworld"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("---")
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		if err != nil {
			log.Fatalf("ListStr get stream err: %v", err)
		}
		// 打印返回值
		log.Println(res.Status)
	}
}
