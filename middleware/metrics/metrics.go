package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metrics"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// Option is metrics option.
type Option func(*options)

// WithRequests with requests counter.
func WithRequests(c metrics.Counter) Option {
	return func(o *options) {
		o.requests = c
	}
}

// WithSeconds with seconds histogram.
func WithSeconds(c metrics.Observer) Option {
	return func(o *options) {
		o.seconds = c
	}
}

type options struct {
	// counter: <client/server>_requests_code_total{kind, full_method, code, reason}
	requests metrics.Counter
	// histogram: <client/server>_requests_seconds_bucket{kind, full_method}
	seconds metrics.Observer
}

// Server is middleware server-side metrics.
func Server(opts ...Option) middleware.Middleware {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				kind      string
				method    string
				code      uint32
				startTime = time.Now()
			)
			if info, ok := transport.FromContext(ctx); ok {
				kind = string(info.Kind)
			}
			if info, ok := middleware.FromContext(ctx); ok {
				method = info.FullMethod
			}
			reply, err := handler(ctx, req)
			if options.requests != nil {
				options.requests.With(kind, method, strconv.Itoa(int(code)), errors.Reason(err)).Inc()
			}
			if options.seconds != nil {
				options.seconds.With(kind, method).Observe(time.Since(startTime).Seconds())
			}
			return reply, err
		}
	}
}

// Client is middleware client-side metrics.
func Client(opts ...Option) middleware.Middleware {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				kind      string
				method    string
				startTime = time.Now()
			)
			if info, ok := transport.FromContext(ctx); ok {
				kind = string(info.Kind)
			}
			if info, ok := middleware.FromContext(ctx); ok {
				method = info.FullMethod
			}
			reply, err := handler(ctx, req)
			if options.requests != nil {
				options.requests.With(kind, method, strconv.Itoa(errors.Code(err)), errors.Reason(err)).Inc()
			}
			if options.seconds != nil {
				options.seconds.With(kind, method).Observe(time.Since(startTime).Seconds())
			}
			return reply, err
		}
	}
}
