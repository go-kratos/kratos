package grpc

import (
	"context"
	"crypto/tls"
	"reflect"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			return handler(ctx, req)
		}
	}
}

func TestUnaryClientInterceptor(t *testing.T) {
	f := unaryClientInterceptor(
		[]middleware.Middleware{EmptyMiddleware()},
		time.Duration(100),
		nil,
		&retry.Strategy{Attempts: 0, Retrier: retry.NewNoRetrier(), Conditions: nil},
	)
	req := &struct{}{}
	resp := &struct{}{}

	err := f(context.TODO(), "hello", req, resp, &grpc.ClientConn{},
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			return nil
		})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWithUnaryInterceptor(t *testing.T) {
	o := &clientOptions{}
	v := []grpc.UnaryClientInterceptor{
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return nil
		},
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
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

func TestDial(t *testing.T) {
	o := &clientOptions{}
	v := []grpc.DialOption{
		grpc.EmptyDialOption{},
	}
	WithOptions(v...)(o)
	if !reflect.DeepEqual(v, o.grpcOpts) {
		t.Errorf("expect %v but got %v", v, o.grpcOpts)
	}
}

func TestDialConn(t *testing.T) {
	_, err := dial(
		context.Background(),
		true,
		WithDiscovery(&mockRegistry{}),
		WithTimeout(10*time.Second),
		WithEndpoint("abc"),
		WithMiddleware(EmptyMiddleware()),
	)
	if err != nil {
		t.Error(err)
	}
}

func TestWithRetryStrategy(t *testing.T) {
	attempt := 0
	testCases := []map[string]interface{}{
		{"method": "test-normal-roundtrip", "expected": 1},
		{"method": "test-retry-count", "expected": 4},
		{"method": "test-error-doing-nothing", "expected": 1},
		{"method": "test-retry-early-exit", "expected": 3},
	}

	f := unaryClientInterceptor(
		[]middleware.Middleware{EmptyMiddleware()},
		time.Duration(100)*time.Second,
		nil,
		&retry.Strategy{
			Attempts:   3,
			Retrier:    retry.NewRetrier(retry.NewConstantBackoff(10*time.Millisecond, 100*time.Millisecond)),
			Conditions: []retry.Condition{retry.NewByCode(13, 14)},
		},
	)
	req := &struct{}{}
	resp := &struct{}{}
	for _, testCase := range testCases {
		_ = f(context.Background(), testCase["method"].(string), req, resp, &grpc.ClientConn{},
			func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				attempt++
				switch method {
				case "test-normal-roundtrip":
					return nil
				case "test-retry-count":
					return status.Error(codes.Internal, "")
				case "test-error-doing-nothing":
					return status.Error(codes.NotFound, "")
				case "test-retry-early-exit":
					if attempt == 3 {
						return nil
					}
					return status.Error(codes.Internal, "")
				}
				t.Error("unexpected error")
				return nil
			})
		if attempt != testCase["expected"].(int) {
			t.Errorf("expected:%v, got:%v", testCase["expected"].(int), attempt)
		}
		attempt = 0
	}
}
