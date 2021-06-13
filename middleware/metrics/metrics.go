package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metrics"
	prom "github.com/go-kratos/kratos/v2/metrics/prometheus"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/prometheus/client_golang/prometheus"
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
	// counter: <client/server>_requests_code_total{kind, operation, code, reason}
	requests metrics.Counter
	// histogram: <client/server>_requests_seconds_bucket{kind, operation}
	seconds metrics.Observer
}

// Server is middleware server-side metrics.
func Server(opts ...Option) middleware.Middleware {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	if options.seconds == nil && options.requests == nil {
		options.seconds, options.requests = defaultPrometheusServer()
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				code      int
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind()
				operation = info.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = int(se.Code)
				reason = se.Reason
			}
			if options.requests != nil {
				options.requests.With(kind, operation, strconv.Itoa(code), reason).Inc()
			}
			if options.seconds != nil {
				options.seconds.With(kind, operation).Observe(time.Since(startTime).Seconds())
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
	if options.seconds == nil && options.requests == nil {
		options.seconds, options.requests = defaultPrometheusClient()
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				code      int
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromClientContext(ctx); ok {
				kind = info.Kind()
				operation = info.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = int(se.Code)
				reason = se.Reason
			}
			if options.requests != nil {
				options.requests.With(kind, operation, strconv.Itoa(code), reason).Inc()
			}
			if options.seconds != nil {
				options.seconds.With(kind, operation).Observe(float64(time.Since(startTime).Milliseconds()))
			}
			return reply, err
		}
	}
}

func defaultPrometheusServer() (metrics.Observer, metrics.Counter) {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "kratos_server",
		Name:      "request_duration_millisecond",
		Help:      "server requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	},
		[]string{"kind", "operation"},
	)
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kratos_server",
			Name:      "api_requests_total",
			Help:      "The total number of processed requests",
		},
		[]string{"kind", "operation", "code", "reason"},
	)
	prometheus.MustRegister(histogram, counter)
	return prom.NewHistogram(histogram), prom.NewCounter(counter)
}

func defaultPrometheusClient() (metrics.Observer, metrics.Counter) {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "kratos_client",
		Name:      "request_duration_millisecond",
		Help:      "client requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	},
		[]string{"kind", "operation"},
	)
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kratos_client",
			Name:      "api_requests_total",
			Help:      "The total number of processed requests",
		},
		[]string{"kind", "operation", "code", "reason"},
	)
	prometheus.MustRegister(histogram, counter)
	return prom.NewHistogram(histogram), prom.NewCounter(counter)
}
