package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go-common/app/service/bbq/push/api/grpc/v1"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"
)

var (
	addr string
)

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
	client := v1.NewPushClient(cc)

	dev := []*v1.Device{
		{
			RegisterID: "1a0018970a8ddcaef46",
			Platform:   1,
			SDK:        1,
			SendNo:     1,
		},
	}

	fmt.Println(client.Notification(context.Background(), &v1.NotificationRequest{
		Dev: dev,
		Body: &v1.NotificationBody{
			Title:   "test title",
			Content: "test content",
			Extra:   "{\"schema\":\"schema\"}",
		},
	}))
}
