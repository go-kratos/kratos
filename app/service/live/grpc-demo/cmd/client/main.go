package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go-common/app/service/live/grpc-demo/api/grpc/v1"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"
)

var name, addr string

func init() {
	flag.StringVar(&name, "name", "lily", "name")
	flag.StringVar(&addr, "addr", "127.0.0.1:9000", "server addr")
}
func main() {
	flag.Parse()
	cfg := &warden.ClientConfig{
		Dial:    xtime.Duration(time.Second * 3),
		Timeout: xtime.Duration(time.Second * 3),
	}
	cc, err := warden.NewClient(cfg).Dial(context.Background(), addr)
	if err != nil {
		log.Fatalf("new client failed!err:=%v", err)
		return
	}
	client := v1.NewGreeterClient(cc)
	resp, err := client.SayHello(context.Background(), &v1.GeeterReq{
		Uid: 123,
	})
	if err != nil {
		log.Fatalf("say hello failed!err:=%v", err)
		return
	}
	fmt.Printf("got HelloReply:%+v", resp)
}
