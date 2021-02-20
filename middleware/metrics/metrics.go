package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metrics"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
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
	// counter: <kind>_<client/server>_requests_code_total{method, path, code}
	requests metrics.Counter
	// histogram: <kind>_<client/server>_requests_seconds_bucket{method, path}
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
				method string
				path   string
				code   int32
			)
			if info, ok := grpc.FromServerContext(ctx); ok {
				method = "POST"
				path = info.FullMethod
			} else if info, ok := http.FromServerContext(ctx); ok {
				method = info.Request.Method
				if route := mux.CurrentRoute(info.Request); route != nil {
					// /path/123 -> /path/{id}
					path, _ = route.GetPathTemplate()
				} else {
					path = info.Request.RequestURI
				}
			}
			startTime := time.Now()
			reply, err := handler(ctx, req)
			if err != nil {
				code = errors.Code(err)
			}
			if options.requests != nil {
				options.requests.With(method, path, strconv.Itoa(int(code))).Inc()
			}
			if options.seconds != nil {
				options.seconds.With(method, path).Observe(time.Since(startTime).Seconds())
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
				method string
				path   string
				code   int32
			)
			if info, ok := grpc.FromClientContext(ctx); ok {
				method = "POST"
				path = info.FullMethod
			} else if info, ok := http.FromClientContext(ctx); ok {
				method = info.Request.Method
				path = info.Request.RequestURI
			}
			startTime := time.Now()
			reply, err := handler(ctx, req)
			if err != nil {
				code = errors.Code(err)
			}
			if options.requests != nil {
				options.requests.With(method, path, strconv.Itoa(int(code))).Inc()
			}
			if options.seconds != nil {
				options.seconds.With(method, path).Observe(time.Since(startTime).Seconds())
			}
			return reply, err
		}
	}
}
