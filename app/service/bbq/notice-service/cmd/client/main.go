package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go-common/app/service/bbq/notice-service/api/v1"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"
)

var (
	addr string
)

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:9003", "server addr")
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
	client := v1.NewNoticeClient(cc)
	req := &v1.NoticeBase{
		Mid:        69121924,
		ActionMid:  111001917,
		SvId:       0,
		Title:      "关注了你",
		Text:       "",
		JumpUrl:    "",
		NoticeType: 3,
		BizType:    3,
		BizId:      0,
		NoticeTime: xtime.Time(time.Now().Unix()),
		Buvid:      "XYDEB30D9F184E3F3EE7536645CB7188E7143",
	}

	resp, err := client.CreateNotice(context.Background(), req)
	fmt.Println(resp, err)
}
