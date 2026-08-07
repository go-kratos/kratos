package metadata

import (
	"context"
	"net/url"
	"strings"
	"unicode"

	"github.com/go-kratos/kratos/v3/metadata"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/transport"
)

// encodeValue percent-encodes v only when necessary: non-ASCII bytes or '%'.
// Printable ASCII values without '%' are written unchanged for backward
// compatibility with older peers that do not percent-decode on read.
func encodeValue(v string) string {
	for _, r := range v {
		if r > unicode.MaxASCII || r == '%' {
			return url.PathEscape(v)
		}
	}
	return v
}

// decodeValue percent-decodes v. Values without '%' are returned as-is.
// Malformed percent-sequences (e.g. a legacy value containing a bare '%')
// fall back to the original string.
func decodeValue(v string) string {
	if !strings.ContainsRune(v, '%') {
		return v
	}
	if decoded, err := url.PathUnescape(v); err == nil {
		return decoded
	}
	return v
}

// Option is metadata option.
type Option func(*options)

type options struct {
	prefix []string
	md     metadata.Metadata
}

func (o *options) hasPrefix(key string) bool {
	k := strings.ToLower(key)
	for _, prefix := range o.prefix {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}

// WithConstants with constant metadata key value.
func WithConstants(md metadata.Metadata) Option {
	return func(o *options) {
		o.md = md
	}
}

// WithPropagatedPrefix with propagated key prefix.
func WithPropagatedPrefix(prefix ...string) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}

// Server is middleware server-side metadata.
func Server(opts ...Option) middleware.Middleware {
	options := &options{
		prefix: []string{"x-md-"}, // x-md-global-, x-md-local
	}
	for _, o := range opts {
		o(options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			md := options.md.Clone()
			header := tr.RequestHeader()
			for _, k := range header.Keys() {
				if options.hasPrefix(k) {
					for _, v := range header.Values(k) {
						md.Add(k, decodeValue(v))
					}
				}
			}
			ctx = metadata.NewServerContext(ctx, md)
			return handler(ctx, req)
		}
	}
}

// Client is middleware client-side metadata.
func Client(opts ...Option) middleware.Middleware {
	options := &options{
		prefix: []string{"x-md-global-"},
	}
	for _, o := range opts {
		o(options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			tr, ok := transport.FromClientContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			header := tr.RequestHeader()
			// x-md-local-
			for k, vList := range options.md {
				for _, v := range vList {
					header.Add(k, encodeValue(v))
				}
			}
			if md, ok := metadata.FromClientContext(ctx); ok {
				for k, vList := range md {
					for _, v := range vList {
						header.Add(k, encodeValue(v))
					}
				}
			}
			// x-md-global-
			if md, ok := metadata.FromServerContext(ctx); ok {
				for k, vList := range md {
					if options.hasPrefix(k) {
						for _, v := range vList {
							header.Add(k, encodeValue(v))
						}
					}
				}
			}
			return handler(ctx, req)
		}
	}
}