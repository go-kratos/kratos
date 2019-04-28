package resolver

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go-common/library/conf/env"
	"go-common/library/naming"
	"go-common/library/net/netutil/breaker"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/rpc/warden/balancer/wrr"
	pb "go-common/library/net/rpc/warden/proto/testproto"
	"go-common/library/net/rpc/warden/resolver"
	xtime "go-common/library/time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var testServerMap map[string]*testServer

func init() {
	testServerMap = make(map[string]*testServer)
}

const testAppID = "main.test"

type testServer struct {
	SayHelloCount int
}

func resetCount() {
	for _, s := range testServerMap {
		s.SayHelloCount = 0
	}
}

func (ts *testServer) SayHello(context.Context, *pb.HelloRequest) (*pb.HelloReply, error) {
	ts.SayHelloCount++
	return &pb.HelloReply{Message: "hello", Success: true}, nil
}

func (ts *testServer) StreamHello(ss pb.Greeter_StreamHelloServer) error {
	panic("not implement error")
}

func createServer(name, listen string) *warden.Server {
	s := warden.NewServer(&warden.ServerConfig{Timeout: xtime.Duration(time.Second)})
	ts := &testServer{}
	testServerMap[name] = ts
	pb.RegisterGreeterServer(s.Server(), ts)
	go func() {
		if err := s.Run(listen); err != nil {
			panic(fmt.Sprintf("run warden server fail! err: %s", err))
		}
	}()
	return s
}

func NSayHello(c pb.GreeterClient, n int) func(*testing.T) {
	return func(t *testing.T) {
		for i := 0; i < n; i++ {
			if _, err := c.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
				t.Fatalf("call sayhello fail! err: %s", err)
			}
		}
	}
}

func createTestClient(t *testing.T) pb.GreeterClient {
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
	conn, err := client.Dial(context.TODO(), "mockdiscovery://authority/main.test")
	if err != nil {
		t.Fatalf("create client fail!err%s", err)
	}
	return pb.NewGreeterClient(conn)
}

var mockResolver *mockDiscoveryBuilder

func newMockDiscoveryBuilder() *mockDiscoveryBuilder {
	return &mockDiscoveryBuilder{
		instances: make(map[string]*naming.Instance),
		watchch:   make(map[string][]*mockDiscoveryResolver),
	}
}
func TestMain(m *testing.M) {
	ctx := context.TODO()
	mockResolver = newMockDiscoveryBuilder()
	resolver.Set(mockResolver)
	s1 := createServer("server1", "127.0.0.1:18081")
	s2 := createServer("server2", "127.0.0.1:18082")
	s3 := createServer("server3", "127.0.0.1:18083")
	s4 := createServer("server4", "127.0.0.1:18084")
	s5 := createServer("server5", "127.0.0.1:18085")
	defer s1.Shutdown(ctx)
	defer s2.Shutdown(ctx)
	defer s3.Shutdown(ctx)
	defer s4.Shutdown(ctx)
	defer s5.Shutdown(ctx)
	os.Exit(m.Run())
}

func TestAddResolver(t *testing.T) {
	mockResolver.registry(testAppID, "server1", "127.0.0.1:18081", map[string]string{})
	c := createTestClient(t)
	t.Run("test_say_hello", NSayHello(c, 10))
	assert.Equal(t, 10, testServerMap["server1"].SayHelloCount)
	resetCount()
}

func TestDeleteResolver(t *testing.T) {
	mockResolver.registry(testAppID, "server1", "127.0.0.1:18081", map[string]string{})
	mockResolver.registry(testAppID, "server2", "127.0.0.1:18082", map[string]string{})
	c := createTestClient(t)
	t.Run("test_say_hello", NSayHello(c, 10))
	assert.Equal(t, 10, testServerMap["server1"].SayHelloCount+testServerMap["server2"].SayHelloCount)

	mockResolver.cancel("server1")
	resetCount()
	time.Sleep(time.Millisecond * 10)
	t.Run("test_say_hello", NSayHello(c, 10))
	assert.Equal(t, 0, testServerMap["server1"].SayHelloCount)

	resetCount()
}

func TestUpdateResolver(t *testing.T) {
	mockResolver.registry(testAppID, "server1", "127.0.0.1:18081", map[string]string{})
	mockResolver.registry(testAppID, "server2", "127.0.0.1:18082", map[string]string{})

	c := createTestClient(t)
	t.Run("test_say_hello", NSayHello(c, 10))
	assert.Equal(t, 10, testServerMap["server1"].SayHelloCount+testServerMap["server2"].SayHelloCount)

	mockResolver.registry(testAppID, "server1", "127.0.0.1:18083", map[string]string{})
	mockResolver.registry(testAppID, "server2", "127.0.0.1:18084", map[string]string{})
	resetCount()
	time.Sleep(time.Millisecond * 10)
	t.Run("test_say_hello", NSayHello(c, 10))
	assert.Equal(t, 0, testServerMap["server1"].SayHelloCount+testServerMap["server2"].SayHelloCount)
	assert.Equal(t, 10, testServerMap["server3"].SayHelloCount+testServerMap["server4"].SayHelloCount)

	resetCount()
}

func TestErrorResolver(t *testing.T) {
	mockResolver := newMockDiscoveryBuilder()
	resolver.Set(mockResolver)
	mockResolver.registry(testAppID, "server1", "127.0.0.1:18081", map[string]string{})
	mockResolver.registry(testAppID, "server6", "127.0.0.1:18086", map[string]string{})

	c := createTestClient(t)
	t.Run("test_say_hello", NSayHello(c, 10))
	assert.Equal(t, 10, testServerMap["server1"].SayHelloCount)

	resetCount()
}

func TestClusterResolver(t *testing.T) {
	mockResolver := newMockDiscoveryBuilder()
	resolver.Set(mockResolver)
	mockResolver.registry(testAppID, "server1", "127.0.0.1:18081", map[string]string{"cluster": "c1"})
	mockResolver.registry(testAppID, "server2", "127.0.0.1:18082", map[string]string{"cluster": "c1"})
	mockResolver.registry(testAppID, "server3", "127.0.0.1:18083", map[string]string{"cluster": "c2"})
	mockResolver.registry(testAppID, "server4", "127.0.0.1:18084", map[string]string{})
	mockResolver.registry(testAppID, "server5", "127.0.0.1:18084", map[string]string{})

	client := warden.NewClient(&warden.ClientConfig{Clusters: []string{"c1"}}, grpc.WithBalancerName(wrr.Name))
	conn, err := client.Dial(context.TODO(), "mockdiscovery://authority/main.test?cluster=c2")
	if err != nil {
		t.Fatalf("create client fail!err%s", err)
	}
	time.Sleep(time.Millisecond * 10)
	cli := pb.NewGreeterClient(conn)
	if _, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
		t.Fatalf("call sayhello fail! err: %s", err)
	}
	if _, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
		t.Fatalf("call sayhello fail! err: %s", err)
	}
	if _, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
		t.Fatalf("call sayhello fail! err: %s", err)
	}
	assert.Equal(t, 1, testServerMap["server1"].SayHelloCount)
	assert.Equal(t, 1, testServerMap["server2"].SayHelloCount)
	assert.Equal(t, 1, testServerMap["server3"].SayHelloCount)

	resetCount()
}

func TestNoClusterResolver(t *testing.T) {
	mockResolver := newMockDiscoveryBuilder()
	resolver.Set(mockResolver)
	mockResolver.registry(testAppID, "server1", "127.0.0.1:18081", map[string]string{"cluster": "c1"})
	mockResolver.registry(testAppID, "server2", "127.0.0.1:18082", map[string]string{"cluster": "c1"})
	mockResolver.registry(testAppID, "server3", "127.0.0.1:18083", map[string]string{"cluster": "c2"})
	mockResolver.registry(testAppID, "server4", "127.0.0.1:18084", map[string]string{})
	client := warden.NewClient(&warden.ClientConfig{}, grpc.WithBalancerName(wrr.Name))
	conn, err := client.Dial(context.TODO(), "mockdiscovery://authority/main.test")
	if err != nil {
		t.Fatalf("create client fail!err%s", err)
	}
	time.Sleep(time.Millisecond * 20)
	cli := pb.NewGreeterClient(conn)
	if _, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
		t.Fatalf("call sayhello fail! err: %s", err)
	}
	if _, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
		t.Fatalf("call sayhello fail! err: %s", err)
	}
	if _, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
		t.Fatalf("call sayhello fail! err: %s", err)
	}
	if _, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
		t.Fatalf("call sayhello fail! err: %s", err)
	}
	assert.Equal(t, 1, testServerMap["server1"].SayHelloCount)
	assert.Equal(t, 1, testServerMap["server2"].SayHelloCount)
	assert.Equal(t, 1, testServerMap["server3"].SayHelloCount)
	assert.Equal(t, 1, testServerMap["server4"].SayHelloCount)

	resetCount()
}

func TestZoneResolver(t *testing.T) {
	mockResolver := newMockDiscoveryBuilder()
	resolver.Set(mockResolver)
	mockResolver.registry(testAppID, "server1", "127.0.0.1:18081", map[string]string{})
	env.Zone = "testsh"
	mockResolver.registry(testAppID, "server2", "127.0.0.1:18082", map[string]string{})
	env.Zone = "hhhh"
	client := warden.NewClient(&warden.ClientConfig{Zone: "testsh"}, grpc.WithBalancerName(wrr.Name))
	conn, err := client.Dial(context.TODO(), "mockdiscovery://authority/main.test")
	if err != nil {
		t.Fatalf("create client fail!err%s", err)
	}
	time.Sleep(time.Millisecond * 10)
	cli := pb.NewGreeterClient(conn)
	if _, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
		t.Fatalf("call sayhello fail! err: %s", err)
	}
	if _, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
		t.Fatalf("call sayhello fail! err: %s", err)
	}
	if _, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
		t.Fatalf("call sayhello fail! err: %s", err)
	}
	assert.Equal(t, 0, testServerMap["server1"].SayHelloCount)
	assert.Equal(t, 3, testServerMap["server2"].SayHelloCount)

	resetCount()
}

func TestSubsetConn(t *testing.T) {
	mockResolver := newMockDiscoveryBuilder()
	resolver.Set(mockResolver)
	mockResolver.registry(testAppID, "server1", "127.0.0.1:18081", map[string]string{})
	mockResolver.registry(testAppID, "server2", "127.0.0.1:18082", map[string]string{})
	mockResolver.registry(testAppID, "server3", "127.0.0.1:18083", map[string]string{})
	mockResolver.registry(testAppID, "server4", "127.0.0.1:18084", map[string]string{})
	mockResolver.registry(testAppID, "server5", "127.0.0.1:18085", map[string]string{})

	client := warden.NewClient(nil)
	conn, err := client.Dial(context.TODO(), "mockdiscovery://authority/main.test?subset=3")
	if err != nil {
		t.Fatalf("create client fail!err%s", err)
	}
	time.Sleep(time.Millisecond * 20)
	cli := pb.NewGreeterClient(conn)
	for i := 0; i < 6; i++ {
		if _, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
			t.Fatalf("call sayhello fail! err: %s", err)
		}
	}
	assert.Equal(t, 2, testServerMap["server2"].SayHelloCount)
	assert.Equal(t, 2, testServerMap["server5"].SayHelloCount)
	assert.Equal(t, 2, testServerMap["server4"].SayHelloCount)
	resetCount()
	mockResolver.cancel("server4")
	time.Sleep(time.Millisecond * 20)
	for i := 0; i < 6; i++ {
		if _, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
			t.Fatalf("call sayhello fail! err: %s", err)
		}
	}
	assert.Equal(t, 2, testServerMap["server5"].SayHelloCount)
	assert.Equal(t, 2, testServerMap["server2"].SayHelloCount)
	assert.Equal(t, 2, testServerMap["server3"].SayHelloCount)
	resetCount()
	mockResolver.registry(testAppID, "server4", "127.0.0.1:18084", map[string]string{})
	time.Sleep(time.Millisecond * 20)
	for i := 0; i < 6; i++ {
		if _, err := cli.SayHello(context.TODO(), &pb.HelloRequest{Age: 1, Name: "hello"}); err != nil {
			t.Fatalf("call sayhello fail! err: %s", err)
		}
	}
	assert.Equal(t, 2, testServerMap["server2"].SayHelloCount)
	assert.Equal(t, 2, testServerMap["server5"].SayHelloCount)
	assert.Equal(t, 2, testServerMap["server4"].SayHelloCount)
}
