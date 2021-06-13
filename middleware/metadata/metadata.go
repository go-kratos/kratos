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
	globalPrefix []string
	md           metadata.Metadata
}

// WithConstants is option with constant metadata key value.
func WithConstants(md metadata.Metadata) Option {
	return func(o *options) {
		o.md = md
	}
}

// GlobalPropagation is option with global propagated key prefix.
func GlobalPropagatedPrefix(prefix []string) Option {
	return func(o *options) {
		o.globalPrefix = append(prefix, prefix...)
	}
}

// Client is middleware client-side metadata.
func Client(opts ...Option) middleware.Middleware {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromClientContext(ctx); ok {
				for k, v := range options.md {
					tr.Header().Set(k, v)
				}
				// passing through the client outgoing metadata
				if cmd, ok := metadata.FromClientContext(ctx); ok {
					for k, v := range cmd {
						tr.Header().Set(k, v)
					}
				}
			}
			return handler(ctx, req)
		}
	}
}

// Server is middleware client-side metadata.
func Server(opts ...Option) middleware.Middleware {
	options := options{
		globalPrefix: []string{"x-md-g-"},
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				var smd metadata.Metadata
				var cmd metadata.Metadata
				for _, k := range tr.Header().Keys() {
					if smd == nil {
						smd = metadata.New()
					}
					val := tr.Header().Get(k)
					smd.Set(k, val)
					for _, prefix := range options.globalPrefix {
						if strings.HasPrefix(strings.ToLower(k), prefix) {
							if cmd == nil {
								cmd = metadata.New()
							}
							cmd.Set(k, val)
						}
					}
				}
				if smd != nil {
					ctx = metadata.NewServerContext(ctx, smd)
				}
				if cmd != nil {
					ctx = metadata.NewClientContext(ctx, cmd)
				}
			}
			return handler(ctx, req)
		}
	}
}
