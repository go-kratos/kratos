package metadata

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// ClientOption is client side metadata option.
type ClientOption func(*options)

// ServerOption is server side metadata option.
type ServerOption func(*options)

type options struct {
	prefix []string
	md     metadata.Metadata
}

// WithConstants with constant metadata key value.
func WithConstants(md metadata.Metadata) ClientOption {
	return func(o *options) {
		o.md = md
	}
}

// WithPropagatedPrefix with propagated key prefix.
func WithPropagatedPrefix(prefix ...string) ServerOption {
	return func(o *options) {
		o.prefix = prefix
	}
}

// Server is middleware client-side metadata.
func Server(opts ...ServerOption) middleware.Middleware {
	options := options{
		prefix: []string{"x-md-"}, // x-md-global-, x-md-local
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				md := metadata.New()
				for _, k := range tr.Header().Keys() {
					for _, prefix := range options.prefix {
						if strings.HasPrefix(strings.ToLower(k), prefix) {
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
func Client(opts ...ClientOption) middleware.Middleware {
	options := options{
		prefix: []string{"x-md-global-"},
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				// x-md-local-
				for k, v := range options.md {
					tr.Header().Set(k, v)
				}
				if md, ok := metadata.FromClientContext(ctx); ok {
					for k, v := range md {
						tr.Header().Set(k, v)
					}
				}
				// x-md-global-
				if md, ok := metadata.FromServerContext(ctx); ok {
					for k, v := range md {
						for _, prefix := range options.prefix {
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
