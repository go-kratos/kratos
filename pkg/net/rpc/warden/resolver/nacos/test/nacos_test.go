package nacos

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/pkg/naming/nacos"
	"github.com/go-kratos/kratos/pkg/net/netutil/breaker"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"os"
	"testing"
	"time"

	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	pb "github.com/go-kratos/kratos/pkg/net/rpc/warden/internal/proto/testproto"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden/resolver"
	xtime "github.com/go-kratos/kratos/pkg/time"
)

type testServer struct {
	name string
}

func (ts *testServer) SayHello(context.Context, *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: ts.name, Success: true}, nil
}

func (ts *testServer) StreamHello(ss pb.Greeter_StreamHelloServer) error {
	panic("not implement error")
}

func createServer(name, listen string) *warden.Server {
	s := warden.NewServer(&warden.ServerConfig{Timeout: xtime.Duration(time.Second)})
	ts := &testServer{name}
	pb.RegisterGreeterServer(s.Server(), ts)
	go func() {
		if err := s.Run(listen); err != nil {
			panic(fmt.Sprintf("run warden server fail! err: %s", err))
		}
	}()
	return s
}

func TestMain(m *testing.M) {
	config := &nacos.Config{
		ServerConfigs: []constant.ServerConfig{
			{
				IpAddr: "192.168.9.102",
				Port:   8848,
			},
		},
		ClientConfig: constant.ClientConfig{
			TimeoutMs:           10 * 1000,
			BeatInterval:        5 * 1000,
			ListenInterval:      30 * 1000,
			NotLoadCacheAtStart: true},
	}
	resolver.Register(nacos.Builder(config))
	ctx := context.TODO()

	s1 := createServer("server1", "127.0.0.1:18001")
	s2 := createServer("server2", "127.0.0.1:18002")
	defer s1.Shutdown(ctx)
	defer s2.Shutdown(ctx)
	os.Exit(m.Run())
}

func createTestClient(t *testing.T, connStr string) pb.GreeterClient {
	client := warden.NewClient(&warden.ClientConfig{
		Dial:    xtime.Duration(time.Second * 10),
		Timeout: xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Window:  xtime.Duration(3 * time.Second),
			Bucket:  10,
			Request: 20,
			K:       1.5,
		},
	})
	conn, err := client.Dial(context.TODO(), connStr)
	if err != nil {
		t.Fatalf("create client fail!err%s", err)
	}
	return pb.NewGreeterClient(conn)
}

func TestNacos(t *testing.T) {
	//cli := createTestClient(t, "discovery://default/127.0.0.1:18003,127.0.0.1:18002")
	cli := createTestClient(t, "nacos://default/server1")
	count := 0
	for i := 0; i < 10; i++ {
		if resp, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
			t.Fatalf("TestNacos: SayHello failed!err:=%v", err)
		} else {
			if resp.Message == "server2" {
				count++
			}
		}
	}
	if count != 10 {
		t.Fatalf("TestNacos: get server2 times must be 10")
	}
}
