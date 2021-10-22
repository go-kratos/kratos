package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	pb "github.com/go-kratos/kratos/v2/internal/testdata/helloworld"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.Name == "error" {
		return nil, errors.BadRequest("custom_error", fmt.Sprintf("invalid argument %s", in.Name))
	}
	if in.Name == "panic" {
		panic("server panic")
	}
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %+v", in.Name)}, nil
}

type testKey struct{}

func TestServer(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, testKey{}, "test")
	srv := NewServer(
		Middleware(
			func(handler middleware.Handler) middleware.Handler {
				return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
					if tr, ok := transport.FromServerContext(ctx); ok {
						if tr.ReplyHeader() != nil {
							tr.ReplyHeader().Set("req_id", "3344")
						}
					}
					return handler(ctx, req)
				}
			}),
		UnaryInterceptor(func(ctx context.Context, req interface{},
			info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return handler(ctx, req)
		}),
		Options(grpc.InitialConnWindowSize(0)),
	)
	pb.RegisterGreeterServer(srv, &server{})

	if e, err := srv.Endpoint(); err != nil || e == nil || strings.HasSuffix(e.Host, ":0") {
		t.Fatal(e, err)
	}

	go func() {
		// start server
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	testClient(t, srv)
	_ = srv.Stop(ctx)
}

func testClient(t *testing.T, srv *Server) {
	u, err := srv.Endpoint()
	if err != nil {
		t.Fatal(err)
	}
	// new a gRPC client
	conn, err := DialInsecure(context.Background(),
		WithEndpoint(u.Host),
		WithOptions(grpc.WithBlock()),
		WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
		WithMiddleware(func(handler middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
				if tr, ok := transport.FromClientContext(ctx); ok {
					header := tr.RequestHeader()
					header.Set("x-md-trace", "2233")
				}
				return handler(ctx, req)
			}
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	client := pb.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "kratos"})
	fmt.Println(err)
	assert.Nil(t, err)
	assert.Equal(t, "Hello kratos", reply.Message)
	_ = conn.Close()
}

func TestNetwork(t *testing.T) {
	o := &Server{}
	v := "abc"
	Network(v)(o)
	assert.Equal(t, v, o.network)
}

func TestAddress(t *testing.T) {
	o := &Server{}
	v := "abc"
	Address(v)(o)
	assert.Equal(t, v, o.address)

	u, err := o.Endpoint()
	assert.NotNil(t, err)
	assert.Nil(t, u)
}

func TestTimeout(t *testing.T) {
	o := &Server{}
	v := time.Duration(123)
	Timeout(v)(o)
	assert.Equal(t, v, o.timeout)
}

func TestMiddleware(t *testing.T) {
	o := &Server{}
	v := []middleware.Middleware{
		func(middleware.Handler) middleware.Handler { return nil },
	}
	Middleware(v...)(o)
	assert.Equal(t, v, o.middleware)
}

type mockLogger struct {
	level log.Level
	key   string
	val   string
}

func (l *mockLogger) Log(level log.Level, keyvals ...interface{}) error {
	l.level = level
	l.key = keyvals[0].(string)
	l.val = keyvals[1].(string)
	return nil
}

func TestLogger(t *testing.T) {
	o := &Server{}
	v := &mockLogger{}
	Logger(v)(o)
	o.log.Log(log.LevelWarn, "foo", "bar")
	assert.Equal(t, "foo", v.key)
	assert.Equal(t, "bar", v.val)
	assert.Equal(t, log.LevelWarn, v.level)
}

func TestTLSConfig(t *testing.T) {
	o := &Server{}
	v := &tls.Config{}
	TLSConfig(v)(o)
	assert.Equal(t, v, o.tlsConf)
}

func TestUnaryInterceptor(t *testing.T) {
	o := &Server{}
	v := []grpc.UnaryServerInterceptor{
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return nil, nil
		},
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return nil, nil
		},
	}
	UnaryInterceptor(v...)(o)
	assert.Equal(t, v, o.ints)
}

func TestOptions(t *testing.T) {
	o := &Server{}
	v := []grpc.ServerOption{
		grpc.EmptyServerOption{},
	}
	Options(v...)(o)
	assert.Equal(t, v, o.grpcOpts)
}

type testResp struct {
	Data string
}

func TestServer_unaryServerInterceptor(t *testing.T) {
	u, err := url.Parse("grpc://hello/world")
	assert.NoError(t, err)
	srv := &Server{
		baseCtx:    context.Background(),
		endpoint:   u,
		middleware: []middleware.Middleware{EmptyMiddleware()},
		timeout:    time.Duration(10),
	}
	req := &struct{}{}
	rv, err := srv.unaryServerInterceptor()(context.TODO(), req, &grpc.UnaryServerInfo{}, func(ctx context.Context, req interface{}) (i interface{}, e error) {
		return &testResp{Data: "hi"}, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "hi", rv.(*testResp).Data)
}
