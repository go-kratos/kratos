package main

import (
	"context"
	"flag"
	"fmt"
	"sync/atomic"
	"time"

	"go-common/library/exp/feature"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
	pb "go-common/library/net/rpc/warden/proto/testproto"
	"go-common/library/net/rpc/warden/resolver"
	"go-common/library/net/rpc/warden/resolver/direct"
)

var addrs string
var cli pb.GreeterClient
var concurrency int
var name string
var req int64
var qps int64

func init() {
	log.Init(&log.Config{Stdout: false})
	flag.StringVar(&addrs, "addr", "127.0.0.1:8000,127.0.0.1:8001", "-addr 127.0.0.1:8080,127.0.0.1:8081")
	flag.IntVar(&concurrency, "c", 3, "-c 5")
	flag.StringVar(&name, "name", "test", "-name test")
}

func main() {
	go calcuQPS()
	feature.DefaultGate.AddFlag(flag.CommandLine)
	flag.Parse()
	feature.DefaultGate.SetFromMap(map[string]bool{"dwrr": true})
	resolver.Register(direct.New())
	c := warden.NewClient(nil)
	conn, err := c.Dial(context.Background(), fmt.Sprintf("direct://d/%s", addrs))
	if err != nil {
		panic(err)
	}
	cli = pb.NewGreeterClient(conn)
	for i := 0; i < concurrency; i++ {
		go func() {
			for {
				say()
				time.Sleep(time.Millisecond * 5)
			}
		}()
	}
	time.Sleep(time.Hour)
}

func calcuQPS() {
	var creq, breq int64
	for {
		time.Sleep(time.Second * 5)
		creq = atomic.LoadInt64(&req)
		delta := creq - breq
		atomic.StoreInt64(&qps, delta/5)
		breq = creq
		fmt.Println("HTTP QPS: ", atomic.LoadInt64(&qps))
	}
}
func say() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	reply, err := cli.SayHello(ctx, &pb.HelloRequest{Name: name, Age: 10})
	if err == nil && reply.Success {
		atomic.AddInt64(&req, 1)
	}
}
