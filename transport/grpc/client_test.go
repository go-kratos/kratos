package grpc

import (
	"context"
	"crypto/tls"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc"
	grpcmd "google.golang.org/grpc/metadata"

	"github.com/go-kratos/kratos/v3/errors"
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

type mockClientStream struct {
	grpc.ClientStream
	ctx context.Context
}

func (m *mockClientStream) Context() context.Context {
	return m.ctx
}

func (m *mockClientStream) SendMsg(any) error { return nil }

func (m *mockClientStream) RecvMsg(any) error { return nil }

type clientStreamCtxKey struct{}

func TestStreamClientInterceptorOnceAndContext(t *testing.T) {
	var calls int
	mw := func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			calls++
			if tr, ok := transport.FromClientContext(ctx); ok {
				tr.RequestHeader().Set("authorization", "Bearer x")
			}
			return next(context.WithValue(ctx, clientStreamCtxKey{}, "user-from-middleware"), req)
		}
	}
	var sawHeader string
	streamer := func(ctx context.Context, _ *grpc.StreamDesc, _ *grpc.ClientConn, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
		md, _ := grpcmd.FromOutgoingContext(ctx)
		if vals := md.Get("authorization"); len(vals) > 0 {
			sawHeader = vals[0]
		}
		return &mockClientStream{ctx: ctx}, nil
	}
	cs, err := streamClientInterceptor([]middleware.Middleware{mw}, nil)(
		context.Background(), &grpc.StreamDesc{}, &grpc.ClientConn{}, "/repro.Stream/Watch", streamer,
	)
	if err != nil {
		t.Fatal(err)
	}
	if err = cs.SendMsg(struct{}{}); err != nil {
		t.Fatal(err)
	}
	if err = cs.RecvMsg(struct{}{}); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Errorf("expect middleware once, got %d", calls)
	}
	got, _ := cs.Context().Value(clientStreamCtxKey{}).(string)
	if got != "user-from-middleware" {
		t.Errorf("expect user-from-middleware, got %q", got)
	}
	if sawHeader != "Bearer x" {
		t.Errorf("expect authorization header, got %q", sawHeader)
	}
}

func TestStreamClientInterceptorRejectsBeforeOpen(t *testing.T) {
	mw := func(middleware.Handler) middleware.Handler {
		return func(context.Context, any) (any, error) {
			return nil, errors.Unauthorized("UNAUTHORIZED", "rejected")
		}
	}
	opened := false
	streamer := func(context.Context, *grpc.StreamDesc, *grpc.ClientConn, string, ...grpc.CallOption) (grpc.ClientStream, error) {
		opened = true
		return &mockClientStream{}, nil
	}
	_, err := streamClientInterceptor([]middleware.Middleware{mw}, nil)(
		context.Background(), &grpc.StreamDesc{}, &grpc.ClientConn{}, "/repro.Stream/Watch", streamer,
	)
	if err == nil {
		t.Fatal("expect middleware error")
	}
	if opened {
		t.Fatal("stream opened after middleware rejected it")
	}
}
