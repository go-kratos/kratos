package warden

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestChainUnaryClient(t *testing.T) {
	var orders []string
	factory := func(name string) grpc.UnaryClientInterceptor {
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			orders = append(orders, name+"-in")
			err := invoker(ctx, method, req, reply, cc, opts...)
			orders = append(orders, name+"-out")
			return err
		}
	}
	handlers := []grpc.UnaryClientInterceptor{factory("h1"), factory("h2"), factory("h3")}
	interceptor := chainUnaryClient(handlers)
	interceptor(context.Background(), "test", nil, nil, nil, func(context.Context, string, interface{}, interface{}, *grpc.ClientConn, ...grpc.CallOption) error {
		return nil
	})
	assert.Equal(t, []string{
		"h1-in",
		"h2-in",
		"h3-in",
		"h3-out",
		"h2-out",
		"h1-out",
	}, orders)
}
