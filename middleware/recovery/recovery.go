package recovery

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

// Latency is recovery latency context key
type Latency struct{}

// ErrUnknownRequest is unknown request error.
var ErrUnknownRequest = errors.InternalServer("UNKNOWN", "unknown request error")

// HandlerFunc is recovery handler func.
type HandlerFunc func(ctx context.Context, req, err any) error

// Option is recovery option.
type Option func(*options)

type options struct {
	handler HandlerFunc
	logger  *slog.Logger
}

// WithHandler with recovery handler.
func WithHandler(h HandlerFunc) Option {
	return func(o *options) {
		o.handler = h
	}
}

// WithLogger sets the logger used to record recovered panics. Defaults to
// [slog.Default].
func WithLogger(logger *slog.Logger) Option {
	return func(o *options) {
		o.logger = logger
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
	logger := op.logger
	if logger == nil {
		logger = slog.Default()
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			startTime := time.Now()
			defer func() {
				if rerr := recover(); rerr != nil {
					buf := make([]byte, 64<<10) //nolint:mnd
					n := runtime.Stack(buf, false)
					buf = buf[:n]
					logger.ErrorContext(ctx, "panic recovered",
						slog.Any("panic", rerr),
						slog.Any("request", req),
						slog.String("stack", string(buf)),
					)
					ctx = context.WithValue(ctx, Latency{}, time.Since(startTime).Seconds())
					err = op.handler(ctx, req, rerr)
				}
			}()
			return handler(ctx, req)
		}
	}
}
