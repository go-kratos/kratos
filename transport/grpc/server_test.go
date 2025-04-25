package grpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/internal/matcher"
	pb "github.com/go-kratos/kratos/v2/internal/testdata/helloworld"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHelloStream(streamServer pb.Greeter_SayHelloStreamServer) error {
	tctx, ok := transport.FromServerContext(streamServer.Context())
	if ok {
		tctx.ReplyHeader().Set("123", "123")
	}
	var cnt uint
	for {
		in, err := streamServer.Recv()
		if err != nil {
			return err
		}
		if in.Name == "error" {
			return errors.BadRequest("custom_error", fmt.Sprintf("invalid argument %s", in.Name))
		}
		if in.Name == "panic" {
			panic("server panic")
		}
		err = streamServer.Send(&pb.HelloReply{
			Message: fmt.Sprintf("hello %s", in.Name),
		})
		if err != nil {
			return err
		}
		cnt++
		if cnt > 1 {
			return nil
		}
	}
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
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
				return func(ctx context.Context, req any) (reply any, err error) {
					if tr, ok := transport.FromServerContext(ctx); ok {
						if tr.ReplyHeader() != nil {
							tr.ReplyHeader().Set("req_id", "3344")
						}
					}
					return handler(ctx, req)
				}
			}),
		UnaryInterceptor(func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
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
		WithUnaryInterceptor(
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				return invoker(ctx, method, req, reply, cc, opts...)
			}),
		WithMiddleware(func(handler middleware.Handler) middleware.Handler {
			return func(ctx context.Context, req any) (reply any, err error) {
				if tr, ok := transport.FromClientContext(ctx); ok {
					header := tr.RequestHeader()
					header.Set("x-md-trace", "2233")
				}
				return handler(ctx, req)
			}
		}),
	)
	defer func() {
		_ = conn.Close()
	}()
	if err != nil {
		t.Fatal(err)
	}
	client := pb.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "kratos"})
	t.Log(err)
	if err != nil {
		t.Errorf("failed to call: %v", err)
	}
	if !reflect.DeepEqual(reply.Message, "Hello kratos") {
		t.Errorf("expect %s, got %s", "Hello kratos", reply.Message)
	}

	streamCli, err := client.SayHelloStream(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		_ = streamCli.CloseSend()
	}()
	err = streamCli.Send(&pb.HelloRequest{Name: "cc"})
	if err != nil {
		t.Error(err)
		return
	}
	reply, err = streamCli.Recv()
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(reply.Message, "hello cc") {
		t.Errorf("expect %s, got %s", "hello cc", reply.Message)
	}
}

func TestNetwork(t *testing.T) {
	o := &Server{}
	v := "abc"
	Network(v)(o)
	if !reflect.DeepEqual(v, o.network) {
		t.Errorf("expect %s, got %s", v, o.network)
	}
}

func TestAddress(t *testing.T) {
	v := "abc"
	o := NewServer(Address(v))
	if !reflect.DeepEqual(v, o.address) {
		t.Errorf("expect %s, got %s", v, o.address)
	}
	u, err := o.Endpoint()
	if err == nil {
		t.Errorf("expect %s, got %s", v, err)
	}
	if u != nil {
		t.Errorf("expect %s, got %s", v, u)
	}
}

func TestTimeout(t *testing.T) {
	o := &Server{}
	v := time.Duration(123)
	Timeout(v)(o)
	if !reflect.DeepEqual(v, o.timeout) {
		t.Errorf("expect %s, got %s", v, o.timeout)
	}
}

func TestTLSConfig(t *testing.T) {
	o := &Server{}
	v := &tls.Config{}
	TLSConfig(v)(o)
	if !reflect.DeepEqual(v, o.tlsConf) {
		t.Errorf("expect %v, got %v", v, o.tlsConf)
	}
}

func TestUnaryInterceptor(t *testing.T) {
	o := &Server{}
	v := []grpc.UnaryServerInterceptor{
		func(context.Context, any, *grpc.UnaryServerInfo, grpc.UnaryHandler) (resp any, err error) {
			return nil, nil
		},
		func(context.Context, any, *grpc.UnaryServerInfo, grpc.UnaryHandler) (resp any, err error) {
			return nil, nil
		},
	}
	UnaryInterceptor(v...)(o)
	if !reflect.DeepEqual(v, o.unaryInts) {
		t.Errorf("expect %v, got %v", v, o.unaryInts)
	}
}

func TestStreamInterceptor(t *testing.T) {
	o := &Server{}
	v := []grpc.StreamServerInterceptor{
		func(any, grpc.ServerStream, *grpc.StreamServerInfo, grpc.StreamHandler) error {
			return nil
		},
		func(any, grpc.ServerStream, *grpc.StreamServerInfo, grpc.StreamHandler) error {
			return nil
		},
	}
	StreamInterceptor(v...)(o)
	if !reflect.DeepEqual(v, o.streamInts) {
		t.Errorf("expect %v, got %v", v, o.streamInts)
	}
}

func TestOptions(t *testing.T) {
	o := &Server{}
	v := []grpc.ServerOption{
		grpc.EmptyServerOption{},
	}
	Options(v...)(o)
	if !reflect.DeepEqual(v, o.grpcOpts) {
		t.Errorf("expect %v, got %v", v, o.grpcOpts)
	}
}

type testResp struct {
	Data string
}

func TestServer_unaryServerInterceptor(t *testing.T) {
	u, err := url.Parse("grpc://hello/world")
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	srv := &Server{
		baseCtx:    context.Background(),
		endpoint:   u,
		timeout:    time.Duration(10),
		middleware: matcher.New(),
	}
	srv.middleware.Use(EmptyMiddleware())
	req := &struct{}{}
	rv, err := srv.unaryServerInterceptor()(context.TODO(), req, &grpc.UnaryServerInfo{}, func(context.Context, any) (any, error) {
		return &testResp{Data: "hi"}, nil
	})
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual("hi", rv.(*testResp).Data) {
		t.Errorf("expect %s, got %s", "hi", rv.(*testResp).Data)
	}
}

type mockServerStream struct {
	ctx      context.Context
	sentMsg  any
	recvMsg  any
	metadata metadata.MD
	grpc.ServerStream
}

func (m *mockServerStream) SetHeader(md metadata.MD) error {
	m.metadata = md
	return nil
}

func (m *mockServerStream) SendHeader(md metadata.MD) error {
	m.metadata = md
	return nil
}

func (m *mockServerStream) SetTrailer(md metadata.MD) {
	m.metadata = md
}

func (m *mockServerStream) Context() context.Context {
	return m.ctx
}

func (m *mockServerStream) SendMsg(msg any) error {
	m.sentMsg = msg
	return nil
}

func (m *mockServerStream) RecvMsg(msg any) error {
	m.recvMsg = msg
	return nil
}

func TestServer_streamServerInterceptor(t *testing.T) {
	u, err := url.Parse("grpc://hello/world")
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	srv := &Server{
		baseCtx:          context.Background(),
		endpoint:         u,
		timeout:          time.Duration(10),
		middleware:       matcher.New(),
		streamMiddleware: matcher.New(),
	}

	srv.streamMiddleware.Use(EmptyMiddleware())

	mockStream := &mockServerStream{
		ctx: srv.baseCtx,
	}

	handler := func(_ any, stream grpc.ServerStream) error {
		resp := &testResp{Data: "stream hi"}
		return stream.SendMsg(resp)
	}

	info := &grpc.StreamServerInfo{
		FullMethod: "/grpc.reflection.v1.ServerReflection/ServerReflectionInfo",
	}

	err = srv.streamServerInterceptor()(nil, mockStream, info, handler)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}

	// Check response
	resp := mockStream.sentMsg.(*testResp)
	if !reflect.DeepEqual("stream hi", resp.Data) {
		t.Errorf("expect %s, got %s", "stream hi", resp.Data)
	}
}

func TestListener(t *testing.T) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	s := &Server{}
	Listener(lis)(s)
	if !reflect.DeepEqual(lis, s.lis) {
		t.Errorf("expect %v, got %v", lis, s.lis)
	}
	if e, err := s.Endpoint(); err != nil || e == nil {
		t.Errorf("expect not empty")
	}
}

func TestStop(t *testing.T) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	tests := []struct {
		name          string
		ctx           context.Context
		cancel        context.CancelFunc
		wantForceStop bool
	}{
		{
			name:          "normal",
			ctx:           context.Background(),
			cancel:        func() {},
			wantForceStop: false,
		},
		{
			name:          "timeout",
			ctx:           timeoutCtx,
			cancel:        cancel,
			wantForceStop: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := net.Listen("tcp", ":0")
			if err != nil {
				t.Fatal(err)
			}
			defer l.Close()

			old := log.GetLogger()
			defer log.SetLogger(old)

			// Create a logger to capture logs
			var logs safeBytesBuffer
			log.SetLogger(log.NewStdLogger(&logs))

			s := NewServer(Listener(l))
			pb.RegisterGreeterServer(s, &server{})

			go func() {
				err := s.Start(context.Background()) //nolint
				if err != nil {
					log.Fatal(err)
				}
			}()

			time.Sleep(100 * time.Millisecond)

			conn, err := DialInsecure(
				context.Background(),
				WithEndpoint(l.Addr().String()),
				WithOptions(grpc.WithBlock()),
			)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()

			go func() {
				client := pb.NewGreeterClient(conn)
				if tt.wantForceStop {
					// Simulate a long-running request
					s, err := client.SayHelloStream(context.Background()) //nolint
					if err != nil {
						log.Fatal(err)
					}
					// Keep the stream open
					for {
						// Intentionally do not send messages, only receive messages
						_, err := s.Recv()
						if err != nil {
							break
						}
					}
				} else {
					_, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "test"}) //nolint
					if err != nil {
						log.Error(err)
					}
				}
			}()

			time.Sleep(100 * time.Millisecond)

			err = s.Stop(tt.ctx)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
				return
			}

			// Check if the stop was forced or graceful
			if tt.wantForceStop {
				if !strings.Contains(logs.String(), "force stop") {
					t.Errorf("Expected force stop\n%s", logs.String())
				}
			} else {
				if strings.Contains(logs.String(), "force stop") {
					t.Errorf("Expected graceful stop\n%s", logs.String())
				}
			}
		})
	}
}

type safeBytesBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *safeBytesBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *safeBytesBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}
