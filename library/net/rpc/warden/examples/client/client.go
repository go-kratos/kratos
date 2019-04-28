package main

import (
	"context"
	"flag"
	"fmt"

	"go-common/library/ecode"
	"go-common/library/net/rpc/warden"
	pb "go-common/library/net/rpc/warden/proto/testproto"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
)

// usage: ./client -grpc.target=test.service=127.0.0.1:8080
func main() {
	flag.Parse()
	conn, err := warden.NewConn("discovery://d/test.service?t=t&cluster=asdasd")
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
