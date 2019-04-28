package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go-common/app/service/bbq/recsys-recall/api/grpc/v1"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"
)

var (
	addr string
)

func test1(client v1.RecsysRecallClient) {
	var infos []*v1.RecallInfo
	// infos = append(infos, &v1.RecallInfo{
	// 	Tag:      "HOT_T1:30",
	// 	Name:     "HOT",
	// 	Priority: 1,
	// 	Limit:    20,
	// })
	// infos = append(infos, &v1.RecallInfo{
	// 	Tag:      "RECALL:HOT_T:10053",
	// 	Name:     "op",
	// 	Scorer:   "default",
	// 	Filter:   "bloomfilter",
	// 	Priority: 2,
	// 	Limit:    10,
	// })
	infos = append(infos, &v1.RecallInfo{
		Tag:      "RECALL:HOT_T:92",
		Name:     "175",
		Scorer:   "default",
		Filter:   "bloomfilter",
		Priority: 1,
		Limit:    5,
	})
	// infos = append(infos, &v1.RecallInfo{
	// 	Tag:   "bbq:recall:tagid:11",
	// 	Name:  "11",
	// 	Limit: 20,
	// })
	// infos = append(infos, &v1.RecallInfo{
	// 	Tag:   "bbq:recall:tagid:802",
	// 	Name:  "802",
	// 	Limit: 20,
	// })
	// infos = append(infos, &v1.RecallInfo{
	// 	Tag:   "bbq:recall:tagid:159",
	// 	Name:  "159",
	// 	Limit: 20,
	// })
	// infos = append(infos, &v1.RecallInfo{
	// 	Tag:      "bbq:recall:tagid:1604",
	// 	Name:     "1604",
	// 	Priority: 20,
	// 	Limit:    20,
	// })
	req := &v1.RecallRequest{
		MID:        5829468,
		BUVID:      "d9972de637d2f3b8939ee628a7ea789b",
		Infos:      infos,
		TotalLimit: 20,
	}
	resp, _ := client.Recall(context.Background(), req)
	fmt.Println(resp)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	for _, v := range resp.List {
		fmt.Println(v)
	}

	// for _, v := range resp.SrcInfo {
	// 	fmt.Println(v)
	// }
}

// func test2(client v1.RecsysRecallClient) {
// 	request := &v1.VideoIndexRequest{
// 		SVIDs: []int64{265375},
// 	}
// 	resp, err := client.VideoIndex(context.Background(), request)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(resp)
// }

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
	client := v1.NewRecsysRecallClient(cc)
	test1(client)
	// test2(client)
}
