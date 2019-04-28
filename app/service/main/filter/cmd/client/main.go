package main

import (
	"context"
	"flag"
	"fmt"
	"go-common/app/service/main/filter/api/grpc/v1"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"
	"log"
	"time"
)

var addr string

func init() {
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
	client := v1.NewFilterClient(cc)
	resp, err := client.Filter(context.Background(), &v1.FilterReq{
		Area:    "reply",
		Message: "习大大",
	})
	if err != nil {
		log.Fatalf("filter failed!err:=%+v", err)
		return
	}
	fmt.Printf("got FilterReply:%+v", resp)
}
