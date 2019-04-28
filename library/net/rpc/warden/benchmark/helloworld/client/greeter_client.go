package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"go-common/library/net/netutil/breaker"
	"go-common/library/net/rpc/warden"
	pb "go-common/library/net/rpc/warden/proto/testproto"
	xtime "go-common/library/time"
)

var (
	ccf = &warden.ClientConfig{
		Dial:    xtime.Duration(time.Second * 10),
		Timeout: xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Window:  xtime.Duration(3 * time.Second),
			Sleep:   xtime.Duration(3 * time.Second),
			Bucket:  10,
			Ratio:   0.3,
			Request: 20,
		},
	}
	cli         pb.GreeterClient
	wg          sync.WaitGroup
	reqSize     int
	concurrency int
	request     int
	all         int64
)

func init() {
	flag.IntVar(&reqSize, "s", 10, "request size")
	flag.IntVar(&concurrency, "c", 10, "concurrency")
	flag.IntVar(&request, "r", 1000, "request per routine")
}

func main() {
	flag.Parse()
	name := randSeq(reqSize)
	cli = newClient()
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go sayHello(&pb.HelloRequest{Name: name})
	}
	wg.Wait()
	fmt.Printf("per request cost %v\n", all/int64(request*concurrency))

}

func sayHello(in *pb.HelloRequest) {
	defer wg.Done()
	now := time.Now()
	for i := 0; i < request; i++ {
		cli.SayHello(context.TODO(), in)
	}
	delta := time.Since(now)
	atomic.AddInt64(&all, int64(delta/time.Millisecond))
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func newClient() (cli pb.GreeterClient) {
	client := warden.NewClient(ccf)
	conn, err := client.Dial(context.TODO(), "127.0.0.1:9999")
	if err != nil {
		return
	}
	cli = pb.NewGreeterClient(conn)
	return
}
