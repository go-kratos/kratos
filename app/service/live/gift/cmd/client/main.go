package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go-common/app/service/live/gift/api/grpc/v1"
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
	client := v1.NewGiftClient(cc)
	//resp, err := client.RoomGiftList(context.Background(), &v1.RoomGiftListReq{
	//	RoomId: 1,
	//	AreaV2ParentId:1,
	//	AreaV2Id:1,
	//	Platform:"android",
	//})

	//resp, err:= client.GiftConfig(context.Background(), &v1.GiftConfigReq{
	//	Platform:"android",
	//	Build:1,
	//})

	resp, err := client.DiscountGiftList(context.Background(), &v1.DiscountGiftListReq{
		//Uid:    88895029,
		Uid:    1,
		Roomid: 1,
		//Ruid: 1,
		AreaV2Id: 89,
	})

	if err != nil {
		log.Fatalf("say hello failed!err:=%v,resp:(%v)", err, resp)
		return
	}
	fmt.Printf("got HelloReply:%+v", resp)
}
