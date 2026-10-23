package metadata

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v3/metadata"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/transport"
)

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

// encodeValue percent-encodes bytes outside the printable ASCII range
// (0x20-0x7E) to make metadata values safe for transport headers.
// gRPC and HTTP require header values to be printable ASCII.
// Pure-ASCII values are returned unchanged for backward compatibility.
func encodeValue(v string) string {
	for i := 0; i < len(v); i++ {
		if v[i] < 0x20 || v[i] > 0x7E {
			var b strings.Builder
			b.Grow(len(v) * 3 / 2)
			for i := 0; i < len(v); i++ {
				c := v[i]
				if c >= 0x20 && c <= 0x7E {
					b.WriteByte(c)
				} else {
					b.WriteString(fmt.Sprintf("%%%02X", c))
				}
			}
			return b.String()
		}
	}
	return v
}

// decodeValue reverses encodeValue, percent-decoding %XX sequences.
func decodeValue(v string) string {
	if !strings.Contains(v, "%") {
		return v
	}
	var b strings.Builder
	b.Grow(len(v))
	for i := 0; i < len(v); i++ {
		if v[i] == '%' && i+2 < len(v) {
			if hi := unhex(v[i+1]); hi >= 0 {
				if lo := unhex(v[i+2]); lo >= 0 {
					b.WriteByte(byte(hi<<4 | lo))
					i += 2
					continue
				}
			}
		}
		b.WriteByte(v[i])
	}
	return b.String()
}

func unhex(c byte) int {
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0')
	case 'a' <= c && c <= 'f':
		return int(c - 'a' + 10)
	case 'A' <= c && c <= 'F':
		return int(c - 'A' + 10)
	}
	return -1
}