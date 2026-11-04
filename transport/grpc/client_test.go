package grpc

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	grpcmd "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/go-kratos/kratos/v3/internal/matcher"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/registry"
	"github.com/go-kratos/kratos/v3/transport"
)

func TestWithEndpoint(t *testing.T) {
	o := &clientOptions{}
	v := "abc"
	WithEndpoint(v)(o)
	if !reflect.DeepEqual(v, o.endpoint) {
		t.Errorf("expect %v but got %v", v, o.endpoint)
	}
}

func TestWithTimeout(t *testing.T) {
	o := &clientOptions{}
	v := time.Duration(123)
	WithTimeout(v)(o)
	if !reflect.DeepEqual(v, o.timeout) {
		t.Errorf("expect %v but got %v", v, o.timeout)
	}
}

func TestWithStreamTimeout(t *testing.T) {
	o := &clientOptions{}
	v := time.Duration(123)
	WithStreamTimeout(v)(o)
	if !reflect.DeepEqual(v, o.streamTimeout) {
		t.Errorf("expect %v but got %v", v, o.streamTimeout)
	}
}

func TestWithMiddleware(t *testing.T) {
	o := &clientOptions{}
	v := []middleware.Middleware{
		func(middleware.Handler) middleware.Handler { return nil },
	}
	WithMiddleware(v...)(o)
	if !reflect.DeepEqual(v, o.middleware) {
		t.Errorf("expect %v but got %v", v, o.middleware)
	}
}

func TestWithStreamMiddleware(t *testing.T) {
	o := &clientOptions{}
	v := []middleware.Middleware{
		func(middleware.Handler) middleware.Handler { return nil },
	}
	WithStreamMiddleware(v...)(o)
	if !reflect.DeepEqual(v, o.streamMiddleware) {
		t.Errorf("expect %v but got %v", v, o.streamInts)
	}
}

type mockRegistry struct{}

func (m *mockRegistry) GetService(_ context.Context, _ string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

func (m *mockRegistry) Watch(_ context.Context, _ string) (registry.Watcher, error) {
	return nil, nil
}

func TestWithDiscovery(t *testing.T) {
	o := &clientOptions{}
	v := &mockRegistry{}
	WithDiscovery(v)(o)
	if !reflect.DeepEqual(v, o.discovery) {
		t.Errorf("expect %v but got %v", v, o.discovery)
	}
}

func TestWithTLSConfig(t *testing.T) {
	o := &clientOptions{}
	v := &tls.Config{}
	WithTLSConfig(v)(o)
	if !reflect.DeepEqual(v, o.tlsConf) {
		t.Errorf("expect %v but got %v", v, o.tlsConf)
	}
}

func EmptyMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			return handler(ctx, req)
		}
	}
}

func TestUnaryClientInterceptor(t *testing.T) {
	f := unaryClientInterceptor([]middleware.Middleware{EmptyMiddleware()}, time.Duration(100), nil)
	req := &struct{}{}
	resp := &struct{}{}

	err := f(context.TODO(), "hello", req, resp, &grpc.ClientConn{},
		func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error {
			return nil
		})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestStreamClientInterceptorTimeout verifies that streamClientInterceptor
// applies the timeout configured with WithStreamTimeout to the stream context,
// and leaves the context unchanged when no stream timeout is set.
func TestStreamClientInterceptorTimeout(t *testing.T) {
	// With a positive timeout the interceptor should apply a deadline.
	interceptorWithTimeout := streamClientInterceptor(nil, 50*time.Millisecond, nil)
	deadlineSet := false
	_, _ = interceptorWithTimeout(
		context.Background(),
		&grpc.StreamDesc{},
		&grpc.ClientConn{},
		"/test.Service/StreamMethod",
		func(ctx context.Context, _ *grpc.StreamDesc, _ *grpc.ClientConn, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
			_, deadlineSet = ctx.Deadline()
			// Return an error so the test does not need a real connection.
			return nil, context.DeadlineExceeded
		},
	)
	if !deadlineSet {
		t.Error("expected a deadline on the stream context when a stream timeout is set, got none")
	}

	// Without a timeout the context should be passed through unchanged.
	interceptorNoTimeout := streamClientInterceptor(nil, 0, nil)
	deadlineSet = false
	_, _ = interceptorNoTimeout(
		context.Background(),
		&grpc.StreamDesc{},
		&grpc.ClientConn{},
		"/test.Service/StreamMethod",
		func(ctx context.Context, _ *grpc.StreamDesc, _ *grpc.ClientConn, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
			_, deadlineSet = ctx.Deadline()
			return nil, context.DeadlineExceeded
		},
	)
	if deadlineSet {
		t.Error("expected no deadline on the stream context when timeout is zero, but deadline was set")
	}
}

type fakeClientStream struct {
	grpc.ClientStream
	err error
}

func (f *fakeClientStream) Header() (grpcmd.MD, error) { return nil, f.err }
func (f *fakeClientStream) SendMsg(any) error          { return f.err }
func (f *fakeClientStream) RecvMsg(any) error          { return f.err }

// TestWrappedClientStreamReleasesTimeout checks that the stream timeout context
// is released when a gRPC client stream ends (any RecvMsg error, a successful
// RecvMsg when the server does not stream, and non-io.EOF errors from SendMsg
// or Header), and kept on other successful calls, on io.EOF from SendMsg or
// Header, and on errors returned only by middleware.
func TestWrappedClientStreamReleasesTimeout(t *testing.T) {
	errInternal := status.Error(codes.Internal, "internal")
	errMiddleware := errors.New("middleware error")
	tests := []struct {
		name          string
		call          string
		err           error
		middlewareFn  bool
		clientStreams bool
		wantReleased  bool
	}{
		{name: "send ok", call: "send", err: nil, wantReleased: false},
		{name: "send io.EOF", call: "send", err: io.EOF, wantReleased: false},
		{name: "send error", call: "send", err: errInternal, wantReleased: true},
		{name: "send middleware error", call: "send", middlewareFn: true, wantReleased: false},
		{name: "recv ok", call: "recv", err: nil, wantReleased: false},
		{name: "recv ok, server does not stream", call: "recv", err: nil, clientStreams: true, wantReleased: true},
		{name: "recv io.EOF", call: "recv", err: io.EOF, wantReleased: true},
		{name: "recv error", call: "recv", err: errInternal, wantReleased: true},
		{name: "recv middleware error", call: "recv", middlewareFn: true, wantReleased: false},
		{name: "header ok", call: "header", err: nil, wantReleased: false},
		{name: "header io.EOF", call: "header", err: io.EOF, wantReleased: false},
		{name: "header error", call: "header", err: errInternal, wantReleased: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := transport.NewClientContext(context.Background(), &Transport{operation: "/test.Service/StreamMethod"})
			ctx, cancel := context.WithTimeout(ctx, time.Hour)
			defer cancel()
			m := matcher.New()
			if tt.middlewareFn {
				m.Use(func(middleware.Handler) middleware.Handler {
					return func(context.Context, any) (any, error) {
						return nil, errMiddleware
					}
				})
			}
			w := &wrappedClientStream{
				ClientStream:  &fakeClientStream{err: tt.err},
				ctx:           ctx,
				middleware:    m,
				cancel:        cancel,
				serverStreams: !tt.clientStreams,
			}
			var err error
			switch tt.call {
			case "send":
				err = w.SendMsg(nil)
			case "recv":
				err = w.RecvMsg(nil)
			case "header":
				_, err = w.Header()
			}
			want := tt.err
			if tt.middlewareFn {
				want = errMiddleware
			}
			if !errors.Is(err, want) {
				t.Fatalf("got error %v, want %v", err, want)
			}
			if released := w.Context().Err() != nil; released != tt.wantReleased {
				t.Errorf("timeout context released = %v, want %v", released, tt.wantReleased)
			}
		})
	}

	t.Run("nil cancel", func(t *testing.T) {
		ctx := transport.NewClientContext(context.Background(), &Transport{operation: "/test.Service/StreamMethod"})
		w := &wrappedClientStream{ClientStream: &fakeClientStream{err: errInternal}, ctx: ctx, middleware: matcher.New()}
		if err := w.SendMsg(nil); !errors.Is(err, errInternal) {
			t.Fatalf("got error %v, want %v", err, errInternal)
		}
		if _, err := w.Header(); !errors.Is(err, errInternal) {
			t.Fatalf("got error %v, want %v", err, errInternal)
		}
		if err := w.RecvMsg(nil); !errors.Is(err, errInternal) {
			t.Fatalf("got error %v, want %v", err, errInternal)
		}
	})
}

// TestStreamSendMsgErrorReleasesTimeout uses a real gRPC stream: a message that
// fails to encode makes SendMsg return a non-io.EOF error, which ends the
// stream, so the stream timeout context must be released too.
func TestStreamSendMsgErrorReleasesTimeout(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	healthpb.RegisterHealthServer(server, health.NewServer())
	go func() { _ = server.Serve(lis) }()
	defer server.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(grpcinsecure.NewCredentials()),
		grpc.WithStreamInterceptor(streamClientInterceptor(nil, time.Hour, nil)),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := conn.NewStream(ctx, &grpc.StreamDesc{ServerStreams: true}, "/grpc.health.v1.Health/Watch")
	if err != nil {
		t.Fatal(err)
	}
	w, ok := stream.(*wrappedClientStream)
	if !ok {
		t.Fatalf("unexpected stream type %T", stream)
	}
	if _, ok := w.Context().Deadline(); !ok {
		t.Fatal("expected a deadline on the stream context")
	}

	// A string field with invalid UTF-8 fails to encode on the client.
	err = stream.SendMsg(&healthpb.HealthCheckRequest{Service: string([]byte{0xff})})
	if status.Code(err) != codes.Internal {
		t.Fatalf("expected codes.Internal from SendMsg, got %v", err)
	}
	if !errors.Is(w.Context().Err(), context.Canceled) {
		t.Errorf("expected the stream timeout context to be canceled, got %v", w.Context().Err())
	}
}

// newStreamTestServer starts a gRPC server whose behavior is chosen by the
// method name, and returns its address.
func newStreamTestServer(t *testing.T) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer(grpc.UnknownServiceHandler(func(_ any, ss grpc.ServerStream) error {
		method, _ := grpc.MethodFromServerStream(ss)
		switch method {
		case "/test.Stream/List": // receive one message, send two
			if err := ss.RecvMsg(&emptypb.Empty{}); err != nil {
				return err
			}
			for i := 0; i < 2; i++ {
				if err := ss.SendMsg(&emptypb.Empty{}); err != nil {
					return err
				}
			}
			return nil
		case "/test.Stream/Collect": // receive until io.EOF, send one
			for {
				err := ss.RecvMsg(&emptypb.Empty{})
				if err == io.EOF {
					return ss.SendMsg(&emptypb.Empty{})
				}
				if err != nil {
					return err
				}
			}
		case "/test.Stream/Slow": // send a message every 50ms for 300ms
			for i := 0; i < 6; i++ {
				time.Sleep(50 * time.Millisecond)
				if err := ss.SendMsg(&emptypb.Empty{}); err != nil {
					return err
				}
			}
			return nil
		case "/test.Stream/Hang":
			<-ss.Context().Done()
			return ss.Context().Err()
		}
		return status.Error(codes.Unimplemented, method)
	}))
	healthpb.RegisterHealthServer(server, health.NewServer())
	go func() { _ = server.Serve(lis) }()
	t.Cleanup(server.Stop)
	return lis.Addr().String()
}

// TestStreamTimeoutOffByDefault checks that WithTimeout does not cut off a
// stream when no stream timeout is configured, while it still applies to
// unary calls.
func TestStreamTimeoutOffByDefault(t *testing.T) {
	addr := newStreamTestServer(t)
	conn, err := NewClient(context.Background(), WithEndpoint(addr), WithTimeout(20*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	stream, err := conn.NewStream(context.Background(), &grpc.StreamDesc{ServerStreams: true}, "/test.Stream/Slow")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := stream.Context().Deadline(); ok {
		t.Fatal("expected no deadline on the stream context without WithStreamTimeout")
	}
	if err = stream.CloseSend(); err != nil {
		t.Fatal(err)
	}
	received := 0
	for {
		err = stream.RecvMsg(&emptypb.Empty{})
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("stream ended after %d messages: %v", received, err)
		}
		received++
	}
	if received != 6 {
		t.Fatalf("received %d messages, want 6", received)
	}

	err = conn.Invoke(context.Background(), "/test.Stream/Hang", &emptypb.Empty{}, &emptypb.Empty{})
	if status.Code(err) != codes.DeadlineExceeded {
		t.Fatalf("expected codes.DeadlineExceeded from a unary call, got %v", err)
	}
}

// TestStreamTimeoutReleasedWhenStreamEnds checks on real streams that the
// stream timeout context is released when the stream ends, and not before.
func TestStreamTimeoutReleasedWhenStreamEnds(t *testing.T) {
	addr := newStreamTestServer(t)
	conn, err := NewClient(context.Background(), WithEndpoint(addr), WithStreamTimeout(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	t.Run("server streaming", func(t *testing.T) {
		stream, err := conn.NewStream(context.Background(), &grpc.StreamDesc{ServerStreams: true}, "/test.Stream/List")
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := stream.Context().Deadline(); !ok {
			t.Fatal("expected a deadline on the stream context")
		}
		if err = stream.SendMsg(&emptypb.Empty{}); err != nil {
			t.Fatal(err)
		}
		if err = stream.CloseSend(); err != nil {
			t.Fatal(err)
		}
		for i := 0; i < 2; i++ {
			if err = stream.RecvMsg(&emptypb.Empty{}); err != nil {
				t.Fatal(err)
			}
			if stream.Context().Err() != nil {
				t.Fatalf("stream timeout context released after message %d, before the stream ended", i+1)
			}
		}
		if err = stream.RecvMsg(&emptypb.Empty{}); err != io.EOF {
			t.Fatalf("expected io.EOF, got %v", err)
		}
		if !errors.Is(stream.Context().Err(), context.Canceled) {
			t.Errorf("expected the stream timeout context to be canceled, got %v", stream.Context().Err())
		}
	})

	t.Run("client streaming with CloseAndRecv", func(t *testing.T) {
		stream, err := conn.NewStream(context.Background(), &grpc.StreamDesc{ClientStreams: true}, "/test.Stream/Collect")
		if err != nil {
			t.Fatal(err)
		}
		client := &grpc.GenericClientStream[emptypb.Empty, emptypb.Empty]{ClientStream: stream}
		for i := 0; i < 3; i++ {
			if err = client.Send(&emptypb.Empty{}); err != nil {
				t.Fatal(err)
			}
		}
		if _, err = client.CloseAndRecv(); err != nil {
			t.Fatal(err)
		}
		if !errors.Is(stream.Context().Err(), context.Canceled) {
			t.Errorf("expected the stream timeout context to be canceled, got %v", stream.Context().Err())
		}
	})
}

func TestWithUnaryInterceptor(t *testing.T) {
	o := &clientOptions{}
	v := []grpc.UnaryClientInterceptor{
		func(context.Context, string, any, any, *grpc.ClientConn, grpc.UnaryInvoker, ...grpc.CallOption) error {
			return nil
		},
		func(context.Context, string, any, any, *grpc.ClientConn, grpc.UnaryInvoker, ...grpc.CallOption) error {
			return nil
		},
	}
	WithUnaryInterceptor(v...)(o)
	if !reflect.DeepEqual(v, o.ints) {
		t.Errorf("expect %v but got %v", v, o.ints)
	}
}

func TestWithOptions(t *testing.T) {
	o := &clientOptions{}
	v := []grpc.DialOption{
		grpc.EmptyDialOption{},
	}
	WithOptions(v...)(o)
	if !reflect.DeepEqual(v, o.grpcOpts) {
		t.Errorf("expect %v but got %v", v, o.grpcOpts)
	}
}

func TestWithHealthCheck(t *testing.T) {
	o := &clientOptions{
		healthCheckConfig: `,"healthCheckConfig":{"serviceName":""}`,
	}
	WithHealthCheck(false)(o)
	if !reflect.DeepEqual("", o.healthCheckConfig) {
		t.Errorf("expect %v but got %v", "", o.healthCheckConfig)
	}
}

func TestNewClientOptions(t *testing.T) {
	o := &clientOptions{}
	v := []grpc.DialOption{
		grpc.EmptyDialOption{},
	}
	WithOptions(v...)(o)
	if !reflect.DeepEqual(v, o.grpcOpts) {
		t.Errorf("expect %v but got %v", v, o.grpcOpts)
	}
}

func TestNewClient(t *testing.T) {
	conn, err := NewClient(
		context.Background(),
		WithDiscovery(&mockRegistry{}),
		WithTimeout(10*time.Second),
		WithEndpoint("abc"),
		WithMiddleware(EmptyMiddleware()),
		WithStreamMiddleware(EmptyMiddleware()),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
}
