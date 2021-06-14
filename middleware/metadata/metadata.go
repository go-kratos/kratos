package metadata

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// Option is metadata option.
type Option func(*options)

type options struct {
	prefix       []string
	globalPrefix []string
	md           metadata.Metadata
}

// WithConstants is option with constant metadata key value.
func WithConstants(md metadata.Metadata) Option {
	return func(o *options) {
		o.md = md
	}
}

// WithGlobalPropagatedPrefix is option with global propagated key prefix.
func WithGlobalPropagatedPrefix(prefix ...string) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}

// PropagatedPrefix is option with global propagated key prefix.
func PropagatedPrefix(prefix ...string) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}

// Server is middleware client-side metadata.
func Server(opts ...Option) middleware.Middleware {
	options := options{
		prefix: []string{"x-md-"},
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				md := metadata.New()
				for _, k := range tr.Header().Keys() {
					key := strings.ToLower(k)
					for _, prefix := range options.prefix {
						if strings.HasPrefix(key, prefix) {
							md.Set(k, tr.Header().Get(k))
							break
						}
					}
				}
				ctx = metadata.NewServerContext(ctx, md)
			}
			return handler(ctx, req)
		}
	}
}

// Client is middleware client-side metadata.
func Client(opts ...Option) middleware.Middleware {
	options := options{
		globalPrefix: []string{"x-md-global-"},
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				for k, v := range options.md {
					tr.Header().Set(k, v)
				}
				if md, ok := metadata.FromClientContext(ctx); ok {
					for k, v := range md {
						tr.Header().Set(k, v)
					}
				}
				if md, ok := metadata.FromServerContext(ctx); ok {
					for k, v := range md {
						for _, prefix := range options.globalPrefix {
							if strings.HasPrefix(k, prefix) {
								tr.Header().Set(k, v)
								break
							}
						}
					}
				}
			}
			return handler(ctx, req)
		}
	}
}
