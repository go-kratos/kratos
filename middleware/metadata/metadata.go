package metadata

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type clientMetadataKey struct{}

func NewClientContext(ctx context.Context, md metadata.Metadata) context.Context {
	return context.WithValue(ctx, clientMetadataKey{}, md)
}

func FromClientContext(ctx context.Context) (metadata.Metadata, bool) {
	md, ok := ctx.Value(clientMetadataKey{}).(metadata.Metadata)
	return md, ok
}

// Option is metadata option.
type Option func(*options)

type options struct {
	prefix string
	md     metadata.Metadata
}

// WithConstants is option with constant metadata key value.
func WithConstants(md metadata.Metadata) Option {
	return func(o *options) {
		o.md = md
	}
}

// WithGlobalPropagation is option with global propagated key prefix.
func WithGlobalPropagation(prefix string) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}

// Client is middleware client-side metadata.
func Client(opts ...Option) middleware.Middleware {
	options := options{
		prefix: "x-md-global-",
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			md := options.md.Clone()
			// passing through the global propagated metadata
			if tr, ok := transport.FromServerContext(ctx); ok {
				for k, v := range tr.Metadata() {
					if strings.HasPrefix(k, options.prefix) {
						md.Set(k, v)
					}
				}
			}
			// passing through the client outgoing metadata
			if cmd, ok := FromClientContext(ctx); ok {
				for k, v := range cmd {
					md.Set(k, v)
				}
			}
			if tr, ok := transport.FromClientContext(ctx); ok {
				tr.WithMetadata(md)
			}
			return handler(ctx, req)
		}
	}
}
