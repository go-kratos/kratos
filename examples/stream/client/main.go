package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/SeeMusic/kratos/examples/stream/hello"
	"github.com/SeeMusic/kratos/v2/middleware/recovery"
	transgrpc "github.com/SeeMusic/kratos/v2/transport/grpc"
)

var wg = sync.WaitGroup{}

func main() {
	conn, err := transgrpc.DialInsecure(
		context.Background(),
		transgrpc.WithEndpoint("127.0.0.1:9001"),
		transgrpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := hello.NewHelloClient(conn)

	wg.Add(3)

	go getNumber(client)
	go uploadLog(client)
	go chat(client)

	wg.Wait()
}

func getNumber(client hello.HelloClient) {
	defer wg.Done()
	stream, err := client.GetNumber(context.Background(), &hello.GetNumberRequest{Data: "2021/08/01"})
	if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ListStr get stream err: %v", err)
		}
		// 打印返回值
		log.Println(res.Number)
	}
}

func uploadLog(client hello.HelloClient) {
	defer wg.Done()
	stream, err := client.UploadLog(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	var number int
	for {
		err := stream.Send(&hello.UploadLogRequest{Log: "log:" + strconv.Itoa(number)})
		if err != nil {
			log.Fatalf("ListStr get stream err: %v", err)
		}
		time.Sleep(time.Millisecond * 50)
		number++
	}
}

func chat(client hello.HelloClient) {
	defer wg.Done()
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	var number int
	for {
		err = stream.Send(&hello.ChatRequest{UpMsg: "kratos:" + strconv.Itoa(number)})
		if err != nil {
			log.Fatalf("ListStr get stream err: %v", err)
		}
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ListStr get stream err: %v", err)
		}
		fmt.Println(res.DownMsg)
		time.Sleep(time.Millisecond * 50)
		number++
	}
}
