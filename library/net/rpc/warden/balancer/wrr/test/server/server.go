package main

import (
	"context"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"math/rand"
	"sync/atomic"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/rpc/warden"
	pb "go-common/library/net/rpc/warden/proto/testproto"
)

var (
	req     int64
	qps     int64
	cpu     int
	errRate int
	sleep   time.Duration
)

func init() {
	log.Init(&log.Config{Stdout: false})
	flag.IntVar(&cpu, "cpu", 3000, "cpu time")
	flag.IntVar(&errRate, "err", 0, "error rate")
	flag.DurationVar(&sleep, "sleep", 0, "sleep time")
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

func main() {
	flag.Parse()
	server := warden.NewServer(nil)
	pb.RegisterGreeterServer(server.Server(), &helloServer{})
	_, err := server.Start()
	if err != nil {
		panic(err)
	}
	calcuQPS()
}

type helloServer struct {
}

func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	atomic.AddInt64(&req, 1)
	if in.Name == "err" {
		if rand.Intn(100) < errRate {
			return nil, ecode.ServiceUnavailable
		}
	}
	time.Sleep(time.Millisecond * time.Duration(in.Age))
	time.Sleep(sleep)
	for i := 0; i < cpu+rand.Intn(cpu); i++ {
		crc32.Checksum([]byte(`testasdwfwfsddsfgwddcscschttp://git.bilibili.co/platform/go-common/merge_requests/new?merge_request%5Bsource_branch%5D=stress%2Fcodel`), crc32.IEEETable)
	}
	return &pb.HelloReply{Message: "Hello " + in.Name, Success: true}, nil
}

func (s *helloServer) StreamHello(ss pb.Greeter_StreamHelloServer) error {
	for i := 0; i < 3; i++ {
		in, err := ss.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		ret := &pb.HelloReply{Message: "Hello " + in.Name, Success: true}
		err = ss.Send(ret)
		if err != nil {
			return err
		}
	}
	return nil
}
