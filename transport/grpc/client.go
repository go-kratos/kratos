package grpc

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/status"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"

	"google.golang.org/grpc"
)

// ClientOption is gRPC client option.
type ClientOption func(o *clientOptions)

// WithEndpoint with client endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// WithTimeout with client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}

// WithMiddleware with client middleware.
func WithMiddleware(m middleware.Middleware) ClientOption {
	return func(o *clientOptions) {
		o.middleware = m
	}
}

// WithRegistry with client registry.
func WithRegistry(r registry.Registry) ClientOption {
	return func(o *clientOptions) {
		o.registry = r
	}
}

// UnaryClientInterceptor with client UnaryClientInterceptor.
func UnaryClientInterceptor(in grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.unaryInt = in
	}
}

// WithOptions with client gRPC options.
func WithOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.grpcOpts = opts
	}
}

// clientOptions is gRPC Client
type clientOptions struct {
	endpoint   string
	timeout    time.Duration
	middleware middleware.Middleware
	registry   registry.Registry
	unaryInt   grpc.UnaryClientInterceptor
	grpcOpts   []grpc.DialOption
}

// Dial returns a GRPC connection.
func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

// DialInsecure returns an insecure GRPC connection.
func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

func dial(ctx context.Context, insecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout: 500 * time.Millisecond,
		middleware: middleware.Chain(
			recovery.Recovery(),
			status.Client(),
		),
	}
	for _, o := range opts {
		o(&options)
	}
	var (
		grpcOpts  = []grpc.DialOption{grpc.WithTimeout(options.timeout)}
		unaryInts = []grpc.UnaryClientInterceptor{unaryClientInterceptor(options.middleware)}
	)
	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}
	if options.registry != nil {
		grpc.WithResolvers(discovery.NewBuilder(options.registry))
	}
	if options.unaryInt != nil {
		unaryInts = append(unaryInts, options.unaryInt)
	}
	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}
	grpcOpts = append(grpcOpts, grpc.WithChainUnaryInterceptor(unaryInts...))
	// creates a client connection to the given endpoint
	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}

// unaryClientInterceptor retruns a unary client interceptor.
func unaryClientInterceptor(m middleware.Middleware) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = transport.NewContext(ctx, transport.Transport{Kind: "gRPC"})
		ctx = NewClientContext(ctx, ClientInfo{FullMethod: method})
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return reply, invoker(ctx, method, req, reply, cc, opts...)
		}
		if m != nil {
			h = m(h)
		}
		_, err := h(ctx, req)
		return err
	}
}
