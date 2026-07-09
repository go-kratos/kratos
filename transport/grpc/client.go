package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcinsecure "google.golang.org/grpc/credentials/insecure"
	grpcmd "google.golang.org/grpc/metadata"

	"github.com/go-kratos/kratos/v3/internal/matcher"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/registry"
	"github.com/go-kratos/kratos/v3/selector"
	"github.com/go-kratos/kratos/v3/selector/wrr"
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/go-kratos/kratos/v3/transport/grpc/resolver/discovery"

	// init resolver
	_ "github.com/go-kratos/kratos/v3/transport/grpc/resolver/direct"
)

func init() {
	if selector.GlobalSelector() == nil {
		selector.SetGlobalSelector(wrr.NewBuilder())
	}
}

// ClientOption is gRPC client option.
type ClientOption func(o *clientOptions)

// WithEndpoint with client endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

// WithSubset with client discovery subset size.
// zero value means subset filter disabled
func WithSubset(size int) ClientOption {
	return func(o *clientOptions) {
		o.subsetSize = size
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

// WithStreamMiddleware with client stream middleware.
func WithStreamMiddleware(m ...middleware.Middleware) ClientOption {
	return func(o *clientOptions) {
		o.streamMiddleware = m
	}
}

// WithDiscovery with client discovery.
func WithDiscovery(d registry.Discovery) ClientOption {
	return func(o *clientOptions) {
		o.discovery = d
	}
}

// WithTLSConfig with TLS config.
func WithTLSConfig(c *tls.Config) ClientOption {
	return func(o *clientOptions) {
		o.tlsConf = c
	}
}

// WithUnaryInterceptor returns a ClientOption that specifies the interceptor for unary RPCs.
func WithUnaryInterceptor(in ...grpc.UnaryClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.ints = in
	}
}

// WithStreamInterceptor returns a ClientOption that specifies the interceptor for streaming RPCs.
func WithStreamInterceptor(in ...grpc.StreamClientInterceptor) ClientOption {
	return func(o *clientOptions) {
		o.streamInts = in
	}
}

// WithOptions with gRPC options.
func WithOptions(opts ...grpc.DialOption) ClientOption {
	return func(o *clientOptions) {
		o.grpcOpts = opts
	}
}

// WithNodeFilter with select filters
func WithNodeFilter(filters ...selector.NodeFilter) ClientOption {
	return func(o *clientOptions) {
		o.filters = filters
	}
}

// WithHealthCheck with health check
func WithHealthCheck(healthCheck bool) ClientOption {
	return func(o *clientOptions) {
		if !healthCheck {
			o.healthCheckConfig = ""
		}
	}
}

// clientOptions is gRPC Client
type clientOptions struct {
	endpoint          string
	subsetSize        int
	tlsConf           *tls.Config
	timeout           time.Duration
	discovery         registry.Discovery
	middleware        []middleware.Middleware
	streamMiddleware  []middleware.Middleware
	ints              []grpc.UnaryClientInterceptor
	streamInts        []grpc.StreamClientInterceptor
	grpcOpts          []grpc.DialOption
	balancerName      string
	filters           []selector.NodeFilter
	healthCheckConfig string
}

// NewClient returns a gRPC client connection.
func NewClient(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	options := clientOptions{
		timeout:           2000 * time.Millisecond,
		balancerName:      balancerName,
		subsetSize:        25,
		healthCheckConfig: `,"healthCheckConfig":{"serviceName":""}`,
	}
	for _, o := range opts {
		o(&options)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	isInsecure := options.tlsConf == nil
	ints := []grpc.UnaryClientInterceptor{
		unaryClientInterceptor(options.middleware, options.timeout, options.filters),
	}
	sints := []grpc.StreamClientInterceptor{
		streamClientInterceptor(options.streamMiddleware, options.timeout, options.filters),
	}

	if len(options.ints) > 0 {
		ints = append(ints, options.ints...)
	}
	if len(options.streamInts) > 0 {
		sints = append(sints, options.streamInts...)
	}
	grpcOpts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]%s}`,
			options.balancerName, options.healthCheckConfig)),
		grpc.WithChainUnaryInterceptor(ints...),
		grpc.WithChainStreamInterceptor(sints...),
	}

	if options.discovery != nil {
		grpcOpts = append(grpcOpts,
			grpc.WithResolvers(
				discovery.NewBuilder(
					options.discovery,
					discovery.WithInsecure(isInsecure),
					discovery.WithTimeout(options.timeout),
					discovery.WithSubset(options.subsetSize),
				)))
	}
	if isInsecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcinsecure.NewCredentials()))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(options.tlsConf)))
	}
	if len(options.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, options.grpcOpts...)
	}
	conn, err := grpc.NewClient(options.endpoint, grpcOpts...)
	if err != nil {
		return nil, err
	}
	conn.Connect()
	return conn, nil
}

func unaryClientInterceptor(ms []middleware.Middleware, timeout time.Duration, filters []selector.NodeFilter) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = transport.NewClientContext(ctx, &Transport{
			endpoint:    cc.Target(),
			operation:   method,
			reqHeader:   headerCarrier{},
			nodeFilters: filters,
		})
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req any) (any, error) {
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
		var p selector.Peer
		ctx = selector.NewPeerContext(ctx, &p)
		_, err := h(ctx, req)
		return err
	}
}

// wrappedClientStream wraps the grpc.ClientStream and applies middleware
type wrappedClientStream struct {
	grpc.ClientStream
	ctx        context.Context
	middleware matcher.Matcher
	// cancel releases the per-stream timeout context created when WithTimeout is
	// configured. It is nil when no timeout is set. Called on the first terminal
	// RecvMsg error (including io.EOF) so the context is not held open past the
	// stream's lifetime.
	cancel context.CancelFunc
}

func (w *wrappedClientStream) Context() context.Context {
	return w.ctx
}

func (w *wrappedClientStream) SendMsg(m any) error {
	h := func(_ context.Context, req any) (any, error) {
		return req, w.ClientStream.SendMsg(m)
	}

	info, ok := transport.FromClientContext(w.ctx)
	if !ok {
		return fmt.Errorf("transport value stored in ctx returns: %v", ok)
	}

	if next := w.middleware.Match(info.Operation()); len(next) > 0 {
		h = middleware.Chain(next...)(h)
	}

	_, err := h(w.ctx, m)
	return err
}

func (w *wrappedClientStream) RecvMsg(m any) error {
	h := func(_ context.Context, req any) (any, error) {
		return req, w.ClientStream.RecvMsg(m)
	}

	info, ok := transport.FromClientContext(w.ctx)
	if !ok {
		return fmt.Errorf("transport value stored in ctx returns: %v", ok)
	}

	if next := w.middleware.Match(info.Operation()); len(next) > 0 {
		h = middleware.Chain(next...)(h)
	}

	_, err := h(w.ctx, m)
	// Release the per-stream timeout context on any terminal error (including
	// io.EOF which signals normal stream completion). This ensures that a
	// context created by streamClientInterceptor with WithTimeout is not held
	// open beyond the stream's actual lifetime.
	if err != nil && w.cancel != nil {
		w.cancel()
	}
	return err
}

func streamClientInterceptor(ms []middleware.Middleware, timeout time.Duration, filters []selector.NodeFilter) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) { // nolint
		ctx = transport.NewClientContext(ctx, &Transport{
			endpoint:    cc.Target(),
			operation:   method,
			reqHeader:   headerCarrier{},
			nodeFilters: filters,
		})
		// Apply the configured per-stream timeout. Unlike unary RPCs where the
		// context is cancelled via defer when the handler returns, streaming RPCs
		// are long-lived: the cancel function is stored in wrappedClientStream and
		// called when RecvMsg receives its first terminal error (including io.EOF),
		// so the context is released exactly when the stream ends.
		var cancel context.CancelFunc
		if timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, timeout)
		}
		var p selector.Peer
		ctx = selector.NewPeerContext(ctx, &p)

		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			if cancel != nil {
				cancel()
			}
			return nil, err
		}

		h := func(_ context.Context, _ any) (any, error) {
			return streamer, nil
		}

		m := matcher.New()
		if len(ms) > 0 {
			m.Use(ms...)
			middleware.Chain(ms...)(h)
		}

		wrappedStream := &wrappedClientStream{
			ClientStream: clientStream,
			ctx:          ctx,
			middleware:   m,
			cancel:       cancel,
		}

		return wrappedStream, nil
	}
}
