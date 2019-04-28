package main

import (
	"context"
	"flag"
	"fmt"
	rpc "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"
	"log"
	"time"
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

	var mid int64 = 6622959 //100011563
	buvid := "123456"
	client := rpc.NewRecsysClient(cc)
	resp, err := client.UpsRecService(context.Background(), &rpc.RecsysRequest{
		MID:    mid,
		BUVID:  buvid,
		Limit:  10,
		Offset: 0,
		SVID:   114888,
	})
	if err != nil {
		log.Fatalf("say hello failed!err:=%v", err)
		return
	}
	fmt.Printf("got Reply: %+v", resp)
}
