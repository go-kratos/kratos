package warden_test

import (
	"context"
	"fmt"
	"io"
	"time"

	"go-common/library/log"
	"go-common/library/net/netutil/breaker"
	"go-common/library/net/rpc/warden"
	pb "go-common/library/net/rpc/warden/proto/testproto"
	xtime "go-common/library/time"

	"google.golang.org/grpc"
)

type helloServer struct {
}

func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
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

func ExampleServer() {
	s := warden.NewServer(&warden.ServerConfig{Timeout: xtime.Duration(time.Second), Addr: ":8080"})
	// apply server interceptor middleware
	s.Use(func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		resp, err := handler(newctx, req)
		return resp, err
	})
	pb.RegisterGreeterServer(s.Server(), &helloServer{})
	s.Start()
}

func ExampleClient() {
	client := warden.NewClient(&warden.ClientConfig{
		Dial:    xtime.Duration(time.Second * 10),
		Timeout: xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Window:  xtime.Duration(3 * time.Second),
			Sleep:   xtime.Duration(3 * time.Second),
			Bucket:  10,
			Ratio:   0.3,
			Request: 20,
		},
	})
	// apply client interceptor middleware
	client.Use(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (ret error) {
		newctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		ret = invoker(newctx, method, req, reply, cc, opts...)
		return
	})
	conn, err := client.Dial(context.Background(), "127.0.0.1:8080")
	if err != nil {
		log.Error("did not connect: %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	name := "2233"
	rp, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name, Age: 18})
	if err != nil {
		log.Error("could not greet: %v", err)
		return
	}
	fmt.Println("rp", *rp)
}
