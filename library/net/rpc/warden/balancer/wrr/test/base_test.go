package test

import (
	"context"
	"io"
	"log"
	"os"
	"sync"
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

	"google.golang.org/grpc"
)

type testBuilder struct {
	addrs []*naming.Instance
}
type testDiscovery struct {
	mu sync.Mutex
	b  *testBuilder
	id string
	ch chan struct{}
}

func (b *testBuilder) Build(id string) naming.Resolver {
	return &testDiscovery{id: id, b: b}
}

func (b *testBuilder) Scheme() string {
	return "testbuilder"
}
func (d *testDiscovery) Fetch(ctx context.Context) (map[string][]*naming.Instance, bool) {
	d.mu.Lock()
	addrs := d.b.addrs
	d.mu.Unlock()
	if len(addrs) == 0 {
		return nil, false
	}
	return map[string][]*naming.Instance{env.Zone: addrs}, true
}

func (d *testDiscovery) Watch() <-chan struct{} {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.ch == nil {
		d.ch = make(chan struct{}, 1)
	}
	return d.ch
}

func (d *testDiscovery) Close() error {
	return nil
}

func (d *testDiscovery) Scheme() string {
	return "discovery"
}

func (d *testDiscovery) set(addrs []*naming.Instance) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.b.addrs = addrs
	select {
	case d.ch <- struct{}{}:
	default:
		return
	}
}

func TestMain(m *testing.M) {
	s1 := runServer(":18080")
	s2 := runServer(":18081")
	s3 := runServer(":18082")
	b = &testBuilder{}
	resolver.Register(b)
	dis = b.Build("test_app").(*testDiscovery)
	go func() {
		time.Sleep(time.Millisecond * 10)
		dis.set([]*naming.Instance{{
			Addrs:    []string{"grpc://127.0.0.1:18080"},
			AppID:    "test_app",
			Metadata: map[string]string{"weight": "100"},
		}, {
			Addrs:    []string{"grpc://127.0.0.1:18081"},
			AppID:    "test_app",
			Metadata: map[string]string{"color": "red"},
		}, {
			Addrs: []string{"grpc://127.0.0.1:18082"},
			AppID: "test_app",
		}})
	}()
	c = newClient()
	time.Sleep(time.Millisecond * 30)
	ret := m.Run()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	s1.Shutdown(ctx)
	s2.Shutdown(ctx)
	s3.Shutdown(ctx)
	os.Exit(ret)
}

type helloServer struct {
	addr string
}

func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: s.addr}, nil
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

func runServer(addr string) *warden.Server {
	server := warden.NewServer(&warden.ServerConfig{Timeout: xtime.Duration(time.Second)})
	pb.RegisterGreeterServer(server.Server(), &helloServer{addr: addr})
	go func() {
		err := server.Run(addr)
		if err != nil {
			panic("run server failed!" + err.Error())
		}
	}()
	return server
}

// NewClient returns a new blank Client instance with a default client interceptor.
// opt can be used to add grpc dial options.
func newClient() (client pb.GreeterClient) {
	c := warden.NewClient(&warden.ClientConfig{
		Dial:    xtime.Duration(time.Second * 10),
		Timeout: xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Window:  xtime.Duration(3 * time.Second),
			Sleep:   xtime.Duration(3 * time.Second),
			Bucket:  10,
			Ratio:   0.3,
			Request: 20,
		},
	},
		grpc.WithBalancerName(wrr.Name),
	)
	conn, err := c.Dial(context.Background(), "discovery://authority/111")
	if err != nil {
		log.Fatalf("can't not connect: %v", err)
	}
	client = pb.NewGreeterClient(conn)
	return
}

var b *testBuilder
var dis *testDiscovery
var c pb.GreeterClient

func TestBalancer(t *testing.T) {
	testBalancerBasic(t)
	testBalancerFailover(t)
	testBalancerUpdateColor(t)
	testBalancerUpdateScore(t)
}

func testBalancerBasic(t *testing.T) {
	time.Sleep(time.Millisecond * 10)
	var idx8080 int
	var idx8082 int
	for i := 0; i < 6; i++ {
		resp, err := c.SayHello(context.Background(), &pb.HelloRequest{Age: 123, Name: "asdasd"})
		if err != nil {
			t.Fatalf("testBalancerBasic: say hello failed!err:=%v", err)
		}
		if resp.Message == ":18082" {
			idx8082++
		} else if resp.Message == ":18080" {
			idx8080++
		}
	}
	if idx8080 != 3 {
		t.Fatalf("testBalancerBasic: server 18080 response times should be 3")
	}
	if idx8082 != 3 {
		t.Fatalf("testBalancerBasic: server 18082 response times should be 3")
	}
}

func testBalancerFailover(t *testing.T) {
	dis.set([]*naming.Instance{{
		Addrs:    []string{"grpc://127.0.0.1:18080"},
		AppID:    "test_app",
		Metadata: map[string]string{"weight": "100"},
	}, {
		Addrs:    []string{"grpc://127.0.0.1:18081"},
		AppID:    "test_app",
		Metadata: map[string]string{"color": "red"},
	}})
	time.Sleep(time.Millisecond * 20)
	var idx8080 int
	var idx8082 int
	for i := 0; i < 4; i++ {
		resp, err := c.SayHello(context.Background(), &pb.HelloRequest{Age: 123, Name: "asdasd"})
		if err != nil {
			t.Fatalf("testBalancerFailover: say hello failed!err:=%v", err)
		}
		if resp.Message == ":18082" {
			idx8082++
		} else if resp.Message == ":18080" {
			idx8080++
		}
	}
	if idx8080 != 4 {
		t.Fatalf("testBalancerFailover: server 8080  response should be 4")
	}
}

func testBalancerUpdateColor(t *testing.T) {
	dis.set([]*naming.Instance{{
		Addrs:    []string{"grpc://127.0.0.1:18080"},
		AppID:    "test_app",
		Metadata: map[string]string{"weight": "100"},
	}, {
		Addrs: []string{"grpc://127.0.0.1:18081"},
		AppID: "test_app",
	}})
	time.Sleep(time.Millisecond * 30)
	var idx8080 int
	var idx8081 int
	for i := 0; i < 4; i++ {
		resp, err := c.SayHello(context.Background(), &pb.HelloRequest{Age: 123, Name: "asdasd"})
		if err != nil {
			t.Fatalf("testBalancerUpdateColor: say hello failed!err:=%v", err)
		}
		if resp.Message == ":18081" {
			idx8081++
		} else if resp.Message == ":18080" {
			idx8080++
		}
	}
	if idx8080 != 2 {
		t.Fatalf("testBalancerUpdateColor: server 8080 response should be 2")
	}
	if idx8081 != 2 {
		t.Fatalf("testBalancerUpdateColor: server 8081 response should be 2")
	}
}

func testBalancerUpdateScore(t *testing.T) {
	dis.set([]*naming.Instance{{
		Addrs:    []string{"grpc://127.0.0.1:18080"},
		AppID:    "test_app",
		Metadata: map[string]string{"weight": "100"},
	}, {
		Addrs:    []string{"grpc://127.0.0.1:18081"},
		AppID:    "test_app",
		Metadata: map[string]string{"weight": "300"},
	}})
	time.Sleep(time.Millisecond * 10)
	var idx8080 int
	var idx8081 int
	for i := 0; i < 4; i++ {
		resp, err := c.SayHello(context.Background(), &pb.HelloRequest{Age: 123, Name: "asdasd"})
		if err != nil {
			t.Fatalf("testBalancerUpdateScore: say hello failed!err:=%v", err)
		}
		if resp.Message == ":18081" {
			idx8081++
		} else if resp.Message == ":18080" {
			idx8080++
		}
	}
	if idx8080 != 1 {
		t.Fatalf("testBalancerUpdateScore: server 8080 response should be 2")
	}
	if idx8081 != 3 {
		t.Fatalf("testBalancerUpdateScore: server 8081 response should be 2")
	}
}
