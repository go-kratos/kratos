package metrics

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
)

// Option is metrics option.
type Option func(*options)

// WithMeter with meter.
func WithMeter(m api.Meter) Option {
	return func(o *options) {
		o.meter = m
	}
}

// WithRequests with requests counter.
func WithRequests(c api.Int64Counter) Option {
	return func(o *options) {
		o.requests = c
	}
}

// WithSeconds with seconds histogram.
func WithSeconds(c api.Float64Histogram) Option {
	return func(o *options) {
		o.seconds = c
	}
}

type options struct {
	meter api.Meter
	// counter: <client/server>_requests_code_total{kind, operation, code, reason}
	requests api.Int64Counter
	// histogram: <client/server>_requests_seconds_bucket{kind, operation}
	seconds api.Float64Histogram
}

// Server is middleware server-side metrics.
func Server(opts ...Option) middleware.Middleware {
	op := options{}
	for _, o := range opts {
		o(&op)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var attributes []attribute.KeyValue
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				attributes = append(attributes,
					attribute.String("kind", info.Kind().String()),
					attribute.String("operation", info.Operation()),
				)
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				attributes = append(attributes,
					attribute.Int("code", int(se.Code)),
					attribute.String("reason", se.Reason),
				)
			}
			if op.requests != nil {
				op.requests.Add(ctx, 1, api.WithAttributes(attributes...))
			}
			if op.seconds != nil {
				op.seconds.Record(ctx, time.Since(startTime).Seconds(), api.WithAttributes(attributes...))
			}
			return reply, err
		}
	}
}

// Client is middleware client-side metrics.
func Client(opts ...Option) middleware.Middleware {
	op := options{}
	for _, o := range opts {
		o(&op)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var attributes []attribute.KeyValue
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				attributes = append(attributes,
					attribute.String("kind", info.Kind().String()),
					attribute.String("operation", info.Operation()),
				)
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				attributes = append(attributes,
					attribute.Int("code", int(se.Code)),
					attribute.String("reason", se.Reason),
				)
			}
			if op.requests != nil {
				op.requests.Add(ctx, 1, api.WithAttributes(attributes...))
			}
			if op.seconds != nil {
				op.seconds.Record(ctx, time.Since(startTime).Seconds(), api.WithAttributes(attributes...))
			}
			return reply, err
		}
	}
}
