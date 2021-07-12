package grpc

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"

	// init resolver
	_ "github.com/go-kratos/kratos/v2/transport/grpc/resolver/direct"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	grpcmd "google.golang.org/grpc/metadata"
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
func WithMiddleware(m ...middleware.Middleware) ClientOption {
	return func(o *clientOptions) {
		o.middleware = m
	}
}

// WithDiscovery with client discovery.
func WithDiscovery(d registry.Discovery) ClientOption {
	return func(o *clientOptions) {
		o.discovery = d
	}
}

// WithUnaryInterceptor returns a DialOption that specifies the interceptor for unary RPCs.
func WithUnaryInterceptor(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.ints = in
	}
}

// WithOptions with gRPC options.
func WithOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.grpcOpts = opts
	}
}

// clientOptions is gRPC Client
type clientOptions struct {
	endpoint   string
	timeout    time.Duration
	discovery  registry.Discovery
	middleware []middleware.Middleware
	ints       []grpc.UnaryClientInterceptor
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
	}
	for _, o := range opts {
		o(&options)
	}
	var ints = []grpc.UnaryClientInterceptor{
		unaryClientInterceptor(options.middleware, options.timeout),
	}
	if len(options.ints) > 0 {
		ints = append(ints, options.ints...)
	}
	var grpcOpts = []grpc.DialOption{
		//todo: grpc.WithBalancerName is deprecated.
		grpc.WithBalancerName(roundrobin.Name), //nolint:staticcheck
		grpc.WithChainUnaryInterceptor(ints...),
	}
	if options.discovery != nil {
		grpcOpts = append(grpcOpts, grpc.WithResolvers(discovery.NewBuilder(options.discovery)))
	}
	if insecure {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}
	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}
	return grpc.DialContext(ctx, options.endpoint, grpcOpts...)
}

func unaryClientInterceptor(ms []middleware.Middleware, timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = transport.NewClientContext(ctx, &Transport{
			endpoint:  cc.Target(),
			operation: method,
			reqHeader: headerCarrier{},
		})
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				header := tr.RequestHeader()
				keys := header.Keys()
				keyvals := make([]string, 0, len(keys))
				for _, k := range keys {
					keyvals = append(keyvals, k, header.Get(k))
				}
				ctx = grpcmd.AppendToOutgoingContext(ctx, keyvals...)
			}
			return reply, invoker(ctx, method, req, reply, cc, opts...)
		}
		if len(ms) > 0 {
			h = middleware.Chain(ms...)(h)
		}
		_, err := h(ctx, req)
		return err
	}
}
