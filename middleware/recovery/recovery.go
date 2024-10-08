package recovery

import (
	"context"
	"runtime"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
)

// Latency is recovery latency context key
type Latency struct{}

// ErrUnknownRequest is unknown request error.
var ErrUnknownRequest = errors.InternalServer("UNKNOWN", "unknown request error")

// HandlerFunc is recovery handler func.
type HandlerFunc func(ctx context.Context, req, err interface{}) error

// Option is recovery option.
type Option func(*options)

type options struct {
	handler HandlerFunc
}

// WithHandler with recovery handler.
func WithHandler(h HandlerFunc) Option {
	return func(o *options) {
		o.handler = h
	}
}

// Recovery is a server middleware that recovers from any panics.
func Recovery(opts ...Option) middleware.Middleware {
	op := options{
		handler: func(context.Context, any, any) error {
			return ErrUnknownRequest
		},
	}
	for _, o := range opts {
		o(&op)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			startTime := time.Now()
			defer func() {
				if rerr := recover(); rerr != nil {
					buf := make([]byte, 64<<10) //nolint:mnd
					n := runtime.Stack(buf, false)
					buf = buf[:n]
					log.Context(ctx).Errorf("%v: %+v\n%s\n", rerr, req, buf)
					ctx = context.WithValue(ctx, Latency{}, time.Since(startTime).Seconds())
					err = op.handler(ctx, req, rerr)
				}
			}()
			return handler(ctx, req)
		}
	}
}
