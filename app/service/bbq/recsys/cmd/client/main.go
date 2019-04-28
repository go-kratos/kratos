package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	rpc "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/library/cache/redis"
	"go-common/library/container/pool"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"
	"log"
	"time"

	"github.com/Dai0522/go-hash/bloomfilter"
)

var name, addr string

func init() {
	flag.StringVar(&name, "name", "lily", "name")
	flag.StringVar(&addr, "addr", "127.0.0.1:9000", "server addr")
}

func bf(svid int64) {
	conf := &redis.Config{
		Config: &pool.Config{
			Active: 10,
			Idle:   10,
		},
		Name:         "recsys-service.user_profile",
		Proto:        "tcp",
		Addr:         "172.16.38.91:6379",
		WriteTimeout: xtime.Duration(1 * time.Second),
		DialTimeout:  xtime.Duration(1 * time.Second),
		ReadTimeout:  xtime.Duration(1 * time.Second),
	}
	rp := redis.NewPool(conf)
	conn := rp.Get(context.Background())
	defer conn.Close()
	b, _ := redis.Bytes(conn.Do("GET", "BBQ:BF:V1:5829468"))
	if b != nil {
		bf, _ := bloomfilter.Load(&b)
		tmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(tmp, uint64(svid))
		fmt.Println(bf.MightContain(tmp))
	}
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

	var MID int64 = 5829468 // 10022647 //100011563
	buvID := "d9972de637d2f3b8939ee628a7ea789b"
	client := rpc.NewRecsysClient(cc)
	resp, err := client.RecService(context.Background(), &rpc.RecsysRequest{
		MID:       MID,
		BUVID:     buvID,
		Limit:     5,
		DebugFlag: true,
		DebugType: "rank",
	})
	if err != nil {
		log.Fatalf("say hello failed!err:=%v", err)
		return
	}
	fmt.Printf("got Reply: %+v", resp)
	for _, v := range resp.List {
		fmt.Println(v.Svid)
		bf(v.Svid)
	}
}
