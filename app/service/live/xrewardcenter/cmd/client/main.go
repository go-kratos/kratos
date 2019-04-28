package main

import (
	"context"
	"flag"
	"log"
	"time"

	"fmt"
	"go-common/app/service/live/xrewardcenter/api/grpc/v1"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"
)

var name, addr string

func init() {
	flag.StringVar(&name, "name", "lily", "name")
	flag.StringVar(&addr, "addr", "127.0.0.1:9004", "server addr")
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

	client := v1.NewAnchorRewardClient(cc)

	resp, err := client.MyReward(context.Background(), &v1.AnchorTaskMyRewardReq{
		Page: 1,
		Uid:  10000,
	})
	fmt.Printf("****** myReward:******* \n %+v \n %v \n", resp, err)

	resp2, err := client.AddReward(context.Background(), &v1.AnchorTaskAddRewardReq{
		RewardId: 1,
		Roomid:   555,
		Source:   1,
		Uid:      10000,
		OrderId:  "test123",
		Lifespan: 1,
	})
	fmt.Printf("*****  addReward:********\n%+v \n %v \n", resp2, err)

	//resp3, err := client.IsViewed(context.Background(), &v1.AnchorTaskIsViewedReq{
	//	Uid: 10000,
	//})
	//fmt.Printf("got IsViewed:%+v", resp3)

	//resp4, _ := client.UseRecord(context.Background(), &v1.AnchorTaskUseRecordReq{
	//	Page: 1,
	//	Uid:  10000,
	//})
	//fmt.Printf("got UseRecord:%+v", resp4)
	//
	//resp5, err := client.UseReward(context.Background(), &v1.AnchorTaskUseRewardReq{
	//	Id:      1,
	//	Uid:     10000,
	//	UsePlat: "ios",
	//})
	//fmt.Printf("got UseReward:%+v", resp5)
}
