package grpc

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestWithEndpoint(t *testing.T) {
	o := &clientOptions{}
	v := "abc"
	WithEndpoint(v)(o)
	assert.Equal(t, v, o.endpoint)
}

func TestWithTimeout(t *testing.T) {
	o := &clientOptions{}
	v := time.Duration(123)
	WithTimeout(v)(o)
	assert.Equal(t, v, o.timeout)
}

func TestWithMiddleware(t *testing.T) {
	o := &clientOptions{}
	v := []middleware.Middleware{
		func(middleware.Handler) middleware.Handler { return nil },
	}
	WithMiddleware(v...)(o)
	assert.Equal(t, v, o.middleware)
}

type mockRegistry struct{}

func (m *mockRegistry) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}
func (m *mockRegistry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return nil, nil
}

func TestWithDiscovery(t *testing.T) {
	o := &clientOptions{}
	v := &mockRegistry{}
	WithDiscovery(v)(o)
	assert.Equal(t, v, o.discovery)
}

func TestWithTLSConfig(t *testing.T) {
	o := &clientOptions{}
	v := &tls.Config{}
	WithTLSConfig(v)(o)
	assert.Equal(t, v, o.tlsConf)
}

func EmptyMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			return handler(ctx, req)
		}
	}
}

func TestUnaryClientInterceptor(t *testing.T) {
	f := unaryClientInterceptor([]middleware.Middleware{EmptyMiddleware()}, time.Duration(100))
	req := &struct{}{}
	resp := &struct{}{}

	err := f(context.TODO(), "hello", req, resp, &grpc.ClientConn{}, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	})
	assert.NoError(t, err)
}

func TestWithUnaryInterceptor(t *testing.T) {
	o := &clientOptions{}
	v := []grpc.UnaryClientInterceptor{
		func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return nil
		},
		func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return nil
		},
	}
	WithUnaryInterceptor(v...)(o)
	assert.Equal(t, v, o.ints)
}

func TestWithOptions(t *testing.T) {
	o := &clientOptions{}
	v := []grpc.DialOption{
		grpc.EmptyDialOption{},
	}
	WithOptions(v...)(o)
	assert.Equal(t, v, o.grpcOpts)
}

func TestDial(t *testing.T) {
	o := &clientOptions{}
	v := []grpc.DialOption{
		grpc.EmptyDialOption{},
	}
	WithOptions(v...)(o)
	assert.Equal(t, v, o.grpcOpts)
}
