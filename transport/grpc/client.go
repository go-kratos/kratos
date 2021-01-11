package grpc

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc"
)

// ClientOption is gRPC client option.
type ClientOption func(o *Client)

// DecodeErrorFunc is encode error func.
type DecodeErrorFunc func(ctx context.Context, err error) error

// Client is grpc transport client.
type Client struct {
	middleware   middleware.Middleware
	errorDecoder DecodeErrorFunc
}

// NewClient new a grpc transport client.
func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		errorDecoder: DefaultErrorDecoder,
	}
	for _, o := range opts {
		o(client)
	}
	return client
}

// Interceptor returns a unary server interceptor.
func (c *Client) Interceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = transport.NewContext(ctx, transport.Transport{Kind: "GRPC"})
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return reply, invoker(ctx, method, req, reply, cc, opts...)
		}
		if c.middleware != nil {
			h = c.middleware(h)
		}
		_, err := h(ctx, req)
		if err != nil {
			return c.errorDecoder(ctx, err)
		}
		return nil
	}
}
