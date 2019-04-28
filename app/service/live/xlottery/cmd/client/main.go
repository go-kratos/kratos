package main

import (
	"context"
	"fmt"

	"go-common/app/service/live/xlottery/api/grpc/v1"
	"go-common/library/net/rpc/warden"
)

//测试grpc
func main() {
	// w := warden.NewClient(nil)
	// conn, err := w.Dial(context.Background(), "localhost:9000")
	// if err != nil {
	// 	panic(err)
	// }
	// client := v1.NewCapsuleClient(conn)
	// resp, err := client.UpdatePoolConfig(context.TODO(), &v1.UpdatePoolConfigReq{Id: 0, CoinTitle: 2, Title: "普通扭蛋币奖池2", StartTime: 1541751430, EndTime: 1541901430, Rule: "普通扭蛋币奖池描述2"})
	// fmt.Println(resp)

	stormTest()
}

func stormTest() {
	w := warden.NewClient(nil)
	conn, err := w.Dial(context.Background(), "172.22.32.242:9000")
	if err != nil {
		panic(err)
	}
	client := v1.NewStormClient(conn)
	resp, err := client.Start(context.TODO(), &v1.StartStormReq{Roomid: 460820, Uid: 88888929, Ruid: 6810576, Num: 1, Beatid: 1, UseShield: true})
	//resp, err := client.CanStart(context.TODO(), &v1.StartStormReq{Roomid: 460820, Uid: 88888929, Ruid: 6810576, Beatid: 1, UseShield: false})
	//resp, err := client.Join(context.TODO(), &v1.JoinStormReq{Id: 14194626115, Roomid: 460820, Mid: 88888929, Platform: "ios"})
	//resp, err := client.Check(context.TODO(), &v1.CheckStormReq{Roomid: 460820, Uid: 88888929})
	//fmt.Println(resp, err)
	fmt.Println(resp, err)
	// client1 := v1.NewCapsuleClient(conn)
	// resp2s, err := client1.DeleteCoin(context.TODO(), &v1.DeleteCoinReq{})
	// fmt.Println(resp2s, err)
}
