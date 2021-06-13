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
	prefix string
	md     metadata.Metadata
}

// WithMetadata with local metadata option.
func WithMetadata(md metadata.Metadata) Option {
	return func(o *options) {
		o.md = md
	}
}

// WithRemotePrefix with remote prefix option.
func WithRemotePrefix(prefix string) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}

// Client is middleware client-side metadata.
func Client(opts ...Option) middleware.Middleware {
	options := options{
		prefix: "x-md-remote-",
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			md := options.md.Clone()
			// passing through the remote metadata
			if tr, ok := transport.FromServerContext(ctx); ok {
				for k, v := range tr.Metadata() {
					if strings.HasPrefix(k, options.prefix) {
						md.Set(k, v)
					}
				}
			}
			if tr, ok := transport.FromClientContext(ctx); ok {
				tr.WithMetadata(md)
			}
			return handler(ctx, req)
		}
	}
}
