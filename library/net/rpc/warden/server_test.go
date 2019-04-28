package warden

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"go-common/library/ecode"
	errpb "go-common/library/ecode/pb"
	"go-common/library/ecode/tip"
	nmd "go-common/library/net/metadata"
	"go-common/library/net/netutil/breaker"
	pb "go-common/library/net/rpc/warden/proto/testproto"
	xtrace "go-common/library/net/trace"
	xtime "go-common/library/time"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	tip.Init(nil)
}

const (
	_separator = "\001"
)

var (
	outPut    []string
	_testOnce sync.Once
	server    *Server

	clientConfig = ClientConfig{
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
	clientConfig2 = ClientConfig{
		Dial:    xtime.Duration(time.Second * 10),
		Timeout: xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Window:  xtime.Duration(3 * time.Second),
			Sleep:   xtime.Duration(3 * time.Second),
			Bucket:  10,
			Ratio:   0.3,
			Request: 20,
		},
		Method: map[string]*ClientConfig{`/testproto.Greeter/SayHello`: {Timeout: xtime.Duration(time.Millisecond * 200)}},
	}
)

type helloServer struct {
	t *testing.T
}

func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.Name == "trace_test" {
		t, isok := xtrace.FromContext(ctx)
		if !isok {
			t = xtrace.New("test title")
			s.t.Fatalf("no trace extracted from server context")
		}
		newCtx := xtrace.NewContext(ctx, t)
		if in.Age == 0 {
			runClient(newCtx, &clientConfig, s.t, "trace_test", 1)
		}
	} else if in.Name == "recovery_test" {
		panic("test recovery")
	} else if in.Name == "graceful_shutdown" {
		time.Sleep(time.Second * 3)
	} else if in.Name == "timeout_test" {
		if in.Age > 10 {
			s.t.Fatalf("can not deliver requests over 10 times because of link timeout")
			return &pb.HelloReply{Message: "Hello " + in.Name, Success: true}, nil
		}
		time.Sleep(time.Millisecond * 10)
		_, err := runClient(ctx, &clientConfig, s.t, "timeout_test", in.Age+1)
		return &pb.HelloReply{Message: "Hello " + in.Name, Success: true}, err
	} else if in.Name == "timeout_test2" {
		if in.Age > 10 {
			s.t.Fatalf("can not deliver requests over 10 times because of link timeout")
			return &pb.HelloReply{Message: "Hello " + in.Name, Success: true}, nil
		}
		time.Sleep(time.Millisecond * 10)
		_, err := runClient(ctx, &clientConfig2, s.t, "timeout_test2", in.Age+1)
		return &pb.HelloReply{Message: "Hello " + in.Name, Success: true}, err
	} else if in.Name == "color_test" {
		if in.Age == 0 {
			resp, err := runClient(ctx, &clientConfig, s.t, "color_test", in.Age+1)
			return resp, err
		}
		color := nmd.String(ctx, nmd.Color)
		return &pb.HelloReply{Message: "Hello " + color, Success: true}, nil
	} else if in.Name == "breaker_test" {
		if rand.Intn(100) <= 50 {
			return nil, status.Errorf(codes.ResourceExhausted, "test")
		}
		return &pb.HelloReply{Message: "Hello " + in.Name, Success: true}, nil
	} else if in.Name == "error_detail" {
		any, _ := ptypes.MarshalAny(&pb.HelloReply{Success: true})
		return nil, &errpb.Error{
			ErrCode:    123456,
			ErrMessage: "test_error_detail",
			ErrDetail:  any,
		}
	} else if in.Name == "ecode_status" {
		reply := &pb.HelloReply{Message: "status", Success: true}
		st, _ := ecode.Error(ecode.RequestErr, "RequestErr").WithDetails(reply)
		return nil, st
	} else if in.Name == "general_error" {
		return nil, fmt.Errorf("haha is error")
	} else if in.Name == "ecode_code_error" {
		return nil, ecode.CreativeArticleTagErr
	} else if in.Name == "pb_error_error" {
		return nil, &errpb.Error{ErrCode: 11122, ErrMessage: "haha"}
	} else if in.Name == "ecode_status_error" {
		return nil, ecode.Error(ecode.RequestErr, "RequestErr")
	} else if in.Name == "test_remote_port" {
		if strconv.Itoa(int(in.Age)) != nmd.String(ctx, nmd.RemotePort) {
			return nil, fmt.Errorf("error port %d", in.Age)
		}
		reply := &pb.HelloReply{Message: "status", Success: true}
		return reply, nil
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

func runServer(t *testing.T, interceptors ...grpc.UnaryServerInterceptor) func() {
	return func() {
		server = NewServer(&ServerConfig{Addr: "127.0.0.1:8080", Timeout: xtime.Duration(time.Second)})
		pb.RegisterGreeterServer(server.Server(), &helloServer{t})
		server.Use(
			func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
				outPut = append(outPut, "1")
				resp, err := handler(ctx, req)
				outPut = append(outPut, "2")
				return resp, err
			},
			func(ctx context.Context, req interface{}, args *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
				outPut = append(outPut, "3")
				resp, err := handler(ctx, req)
				outPut = append(outPut, "4")
				return resp, err
			})
		if _, err := server.Start(); err != nil {
			t.Fatal(err)
		}
	}
}

func runClient(ctx context.Context, cc *ClientConfig, t *testing.T, name string, age int32, interceptors ...grpc.UnaryClientInterceptor) (resp *pb.HelloReply, err error) {
	client := NewClient(cc)
	client.Use(interceptors...)
	conn, err := client.Dial(context.Background(), "127.0.0.1:8080")
	if err != nil {
		panic(fmt.Errorf("did not connect: %v,req: %v %v", err, name, age))
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	resp, err = c.SayHello(ctx, &pb.HelloRequest{Name: name, Age: age})
	return
}

func Test_Warden(t *testing.T) {
	go func() {
		time.Sleep(time.Second * 10)
		panic("run test warden timeout,exit now!")
	}()
	xtrace.Init(&xtrace.Config{Addr: "127.0.0.1:9982", Proto: "udp", Timeout: xtime.Duration(time.Second * 3)})
	go _testOnce.Do(runServer(t))
	go runClient(context.Background(), &clientConfig, t, "trace_test", 0)
	testTrace(t, 9982, false)
	testInterceptorChain(t)
	testValidation(t)
	testServerRecovery(t)
	testClientRecovery(t)
	testErrorDetail(t)
	testECodeStatus(t)
	testColorPass(t)
	testRemotePort(t)
	testLinkTimeout(t)
	testClientConfig(t)
	testBreaker(t)
	testAllErrorCase(t)
	testGracefulShutDown(t)
}

func testValidation(t *testing.T) {
	_, err := runClient(context.Background(), &clientConfig, t, "", 0)
	if !ecode.RequestErr.Equal(err) {
		t.Fatalf("testValidation should return ecode.RequestErr,but is %v", err)
	}
}

func testAllErrorCase(t *testing.T) {
	// } else if in.Name == "general_error" {
	// 	return nil, fmt.Errorf("haha is error")
	// } else if in.Name == "ecode_code_error" {
	// 	return nil, ecode.CreativeArticleTagErr
	// } else if in.Name == "pb_error_error" {
	// 	return nil, &errpb.Error{ErrCode: 11122, ErrMessage: "haha"}
	// } else if in.Name == "ecode_status_error" {
	// 	return nil, ecode.Error(ecode.RequestErr, "RequestErr")
	// }
	ctx := context.Background()
	t.Run("general_error", func(t *testing.T) {
		_, err := runClient(ctx, &clientConfig, t, "general_error", 0)
		assert.Contains(t, err.Error(), "haha")
		ec := ecode.Cause(err)
		assert.Equal(t, -500, ec.Code())
		// remove this assert in future
		assert.Equal(t, "服务器错误", ec.Message())
	})
	t.Run("ecode_code_error", func(t *testing.T) {
		_, err := runClient(ctx, &clientConfig, t, "ecode_code_error", 0)
		ec := ecode.Cause(err)
		assert.Equal(t, ecode.CreativeArticleTagErr.Code(), ec.Code())
		// remove this assert in future
		assert.Equal(t, "标签错误", ec.Message())
	})
	t.Run("pb_error_error", func(t *testing.T) {
		_, err := runClient(ctx, &clientConfig, t, "pb_error_error", 0)
		ec := ecode.Cause(err)
		assert.Equal(t, 11122, ec.Code())
		assert.Equal(t, "haha", ec.Message())
	})
	t.Run("ecode_status_error", func(t *testing.T) {
		_, err := runClient(ctx, &clientConfig, t, "ecode_status_error", 0)
		ec := ecode.Cause(err)
		assert.Equal(t, ecode.RequestErr.Code(), ec.Code())
		assert.Equal(t, "RequestErr", ec.Message())
	})
}

func testBreaker(t *testing.T) {
	client := NewClient(&clientConfig)
	conn, err := client.Dial(context.Background(), "127.0.0.1:8080")
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	for i := 0; i < 35; i++ {
		_, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "breaker_test"})
		if err != nil {
			if ecode.ServiceUnavailable.Equal(err) {
				return
			}
		}
	}
	t.Fatalf("testBreaker failed!No breaker was triggered")
}

func testColorPass(t *testing.T) {
	ctx := nmd.NewContext(context.Background(), nmd.MD{
		nmd.Color: "red",
	})
	resp, err := runClient(ctx, &clientConfig, t, "color_test", 0)
	if err != nil {
		t.Fatalf("testColorPass  return error %v", err)
	}
	if resp == nil || resp.Message != "Hello red" {
		t.Fatalf("testColorPass resp.Message must be red,%v", *resp)
	}
}

func testRemotePort(t *testing.T) {
	ctx := nmd.NewContext(context.Background(), nmd.MD{
		nmd.RemotePort: "8000",
	})
	_, err := runClient(ctx, &clientConfig, t, "test_remote_port", 8000)
	if err != nil {
		t.Fatalf("testRemotePort return error %v", err)
	}
}

func testLinkTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()
	_, err := runClient(ctx, &clientConfig, t, "timeout_test", 0)
	if err == nil {
		t.Fatalf("testLinkTimeout must return error")
	}
	if !ecode.Deadline.Equal(err) {
		t.Fatalf("testLinkTimeout must return error RPCDeadline,err:%v", err)
	}

}
func testClientConfig(t *testing.T) {
	_, err := runClient(context.Background(), &clientConfig2, t, "timeout_test2", 0)
	if err == nil {
		t.Fatalf("testLinkTimeout must return error")
	}
	if !ecode.Deadline.Equal(err) {
		t.Fatalf("testLinkTimeout must return error RPCDeadline,err:%v", err)
	}
}

func testGracefulShutDown(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := runClient(context.Background(), &clientConfig, t, "graceful_shutdown", 0)
			if err != nil {
				panic(fmt.Errorf("run graceful_shutdown client return(%v)", err))
			}
			if !resp.Success || resp.Message != "Hello graceful_shutdown" {
				panic(fmt.Errorf("run graceful_shutdown client return(%v,%v)", err, *resp))
			}
		}()
	}
	go func() {
		time.Sleep(time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		server.Shutdown(ctx)
	}()
	wg.Wait()
}

func testClientRecovery(t *testing.T) {
	ctx := context.Background()
	client := NewClient(&clientConfig)
	client.Use(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (ret error) {
		invoker(ctx, method, req, reply, cc, opts...)
		panic("client recovery test")
	})

	conn, err := client.Dial(ctx, "127.0.0.1:8080")
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	_, err = c.SayHello(ctx, &pb.HelloRequest{Name: "other_test", Age: 0})
	if err == nil {
		t.Fatalf("recovery must return error")
	}
	e, ok := errors.Cause(err).(ecode.Codes)
	if !ok {
		t.Fatalf("recovery must return ecode error")
	}

	if !ecode.ServerErr.Equal(e) {
		t.Fatalf("recovery must return ecode.RPCClientErr")
	}
}

func testServerRecovery(t *testing.T) {
	ctx := context.Background()
	client := NewClient(&clientConfig)

	conn, err := client.Dial(ctx, "127.0.0.1:8080")
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	_, err = c.SayHello(ctx, &pb.HelloRequest{Name: "recovery_test", Age: 0})
	if err == nil {
		t.Fatalf("recovery must return error")
	}
	e, ok := errors.Cause(err).(ecode.Codes)
	if !ok {
		t.Fatalf("recovery must return ecode error")
	}

	if e.Code() != ecode.ServerErr.Code() {
		t.Fatalf("recovery must return ecode.ServerErr")
	}
}

func testInterceptorChain(t *testing.T) {
	// NOTE: don't delete this sleep
	time.Sleep(time.Millisecond)
	if outPut[0] != "1" || outPut[1] != "3" || outPut[2] != "1" || outPut[3] != "3" || outPut[4] != "4" || outPut[5] != "2" || outPut[6] != "4" || outPut[7] != "2" {
		t.Fatalf("outPut shoud be [1 3 1 3 4 2 4 2]!")
	}
}

func testErrorDetail(t *testing.T) {
	_, err := runClient(context.Background(), &clientConfig2, t, "error_detail", 0)
	if err == nil {
		t.Fatalf("testErrorDetail must return error")
	}
	if ec, ok := errors.Cause(err).(ecode.Codes); !ok {
		t.Fatalf("testErrorDetail must return ecode error")
	} else if ec.Code() != 123456 || ec.Message() != "test_error_detail" || len(ec.Details()) == 0 {
		t.Fatalf("testErrorDetail must return code:123456 and message:test_error_detail, code: %d, message: %s, details length: %d", ec.Code(), ec.Message(), len(ec.Details()))
	} else if _, ok := ec.Details()[0].(*pb.HelloReply); !ok {
		t.Fatalf("expect get pb.HelloReply")
	}
}

func testECodeStatus(t *testing.T) {
	_, err := runClient(context.Background(), &clientConfig2, t, "ecode_status", 0)
	if err == nil {
		t.Fatalf("testECodeStatus must return error")
	}
	st, ok := errors.Cause(err).(*ecode.Status)
	if !ok {
		t.Fatalf("testECodeStatus must return *ecode.Status")
	}
	if st.Code() != int(ecode.RequestErr) && st.Message() != "RequestErr" {
		t.Fatalf("testECodeStatus must return code: -400, message: RequestErr get: code: %d, message: %s", st.Code(), st.Message())
	}
	detail := st.Details()[0].(*pb.HelloReply)
	if !detail.Success || detail.Message != "status" {
		t.Fatalf("wrong detail")
	}
}

func testTrace(t *testing.T, port int, isStream bool) {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: port})
	if err != nil {
		t.Fatalf("listent udp failed, %v", err)
		return
	}
	data := make([]byte, 1024)
	strs := make([][]string, 0)
	for {
		var n int
		n, _, err = listener.ReadFromUDP(data)
		if err != nil {
			t.Fatalf("read from udp faild, %v", err)
		}
		str := strings.Split(string(data[:n]), _separator)
		strs = append(strs, str)

		if len(strs) == 2 {
			break
		}
	}
	if len(strs[0]) == 0 || len(strs[1]) == 0 {
		t.Fatalf("trace str's length must be greater than 0")
	}

}

func BenchmarkServer(b *testing.B) {
	server := NewServer(&ServerConfig{Addr: "127.0.0.1:8080", Timeout: xtime.Duration(time.Second)})
	go func() {
		pb.RegisterGreeterServer(server.Server(), &helloServer{})
		if _, err := server.Start(); err != nil {
			os.Exit(0)
			return
		}
	}()
	defer func() {
		server.Server().Stop()
	}()
	client := NewClient(&clientConfig)
	conn, err := client.Dial(context.Background(), "127.0.0.1:8080")
	if err != nil {
		conn.Close()
		b.Fatalf("did not connect: %v", err)
	}
	b.ResetTimer()
	b.RunParallel(func(parab *testing.PB) {
		for parab.Next() {
			c := pb.NewGreeterClient(conn)
			resp, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "benchmark_test", Age: 1})
			if err != nil {
				conn.Close()
				b.Fatalf("c.SayHello failed: %v,req: %v %v", err, "benchmark", 1)
			}
			if !resp.Success {
				b.Error("repsonse not success!")
			}
		}
	})
	conn.Close()
}

func TestParseDSN(t *testing.T) {
	dsn := "tcp://0.0.0.0:80/?timeout=100ms&idleTimeout=120s&keepaliveInterval=120s&keepaliveTimeout=20s&maxLife=4h&closeWait=3s"
	config := parseDSN(dsn)
	if config.Network != "tcp" || config.Addr != "0.0.0.0:80" || time.Duration(config.Timeout) != time.Millisecond*100 ||
		time.Duration(config.IdleTimeout) != time.Second*120 || time.Duration(config.KeepAliveInterval) != time.Second*120 ||
		time.Duration(config.MaxLifeTime) != time.Hour*4 || time.Duration(config.ForceCloseWait) != time.Second*3 || time.Duration(config.KeepAliveTimeout) != time.Second*20 {
		t.Fatalf("parseDSN(%s) not compare config result(%+v)", dsn, config)
	}

	dsn = "unix:///temp/warden.sock?timeout=300ms"
	config = parseDSN(dsn)
	if config.Network != "unix" || config.Addr != "/temp/warden.sock" || time.Duration(config.Timeout) != time.Millisecond*300 {
		t.Fatalf("parseDSN(%s) not compare config result(%+v)", dsn, config)
	}
}
