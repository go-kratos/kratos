package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/bilibili/Kratos/pkg/ecode"
	"github.com/bilibili/Kratos/pkg/log"
	"github.com/bilibili/Kratos/pkg/net/rpc/warden"
	pb "github.com/bilibili/Kratos/pkg/net/rpc/warden/internal/proto/testproto"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
)

// usage: ./client -grpc.target=test.service=127.0.0.1:8080
func main() {
	log.Init(&log.Config{Stdout: true})
	flag.Parse()
	conn, err := warden.NewConn("127.0.0.1:8081")
	if err != nil {
		panic(err)
	}
	cli := pb.NewGreeterClient(conn)
	normalCall(cli)
	errDetailCall(cli)
}

func normalCall(cli pb.GreeterClient) {
	reply, err := cli.SayHello(context.Background(), &pb.HelloRequest{Name: "tom", Age: 23})
	if err != nil {
		panic(err)
	}
	fmt.Println("get reply:", *reply)
}

func errDetailCall(cli pb.GreeterClient) {
	_, err := cli.SayHello(context.Background(), &pb.HelloRequest{Name: "err_detail_test", Age: 12})
	if err != nil {
		any := errors.Cause(err).(ecode.Codes).Details()[0].(*any.Any)
		var re pb.HelloReply
		err := ptypes.UnmarshalAny(any, &re)
		if err == nil {
			fmt.Printf("cli.SayHello get error detail!detail:=%v", re)
		}
		return
	}
}
