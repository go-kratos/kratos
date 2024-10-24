package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc/codes"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http/status"
)

const (
	metricLabelKind      = "kind"
	metricLabelOperation = "operation"
	metricLabelCode      = "code"
	metricLabelReason    = "reason"
)

const (
	DefaultServerSecondsHistogramName = "server_requests_seconds_bucket"
	DefaultServerRequestsCounterName  = "server_requests_code_total"
	DefaultClientSecondsHistogramName = "client_requests_seconds_bucket"
	DefaultClientRequestsCounterName  = "client_requests_code_total"
)

// Option is metrics option.
type Option func(*options)

// WithRequests with requests counter.
func WithRequests(c metric.Int64Counter) Option {
	return func(o *options) {
		o.requests = c
	}
}

// WithSeconds with seconds histogram.
// notice: the record unit in current middleware is s(Seconds)
func WithSeconds(histogram metric.Float64Histogram) Option {
	return func(o *options) {
		o.seconds = histogram
	}
}

// DefaultRequestsCounter
// return metric.Int64Counter for WithRequests
// suggest histogramName = <client/server>_requests_code_total
func DefaultRequestsCounter(meter metric.Meter, histogramName string) (metric.Int64Counter, error) {
	return meter.Int64Counter(histogramName, metric.WithUnit("{call}"))
}

// DefaultSecondsHistogram
// return metric.Float64Histogram for WithSeconds
// suggest histogramName = <client/server>_requests_seconds_bucket
func DefaultSecondsHistogram(meter metric.Meter, histogramName string) (metric.Float64Histogram, error) {
	return meter.Float64Histogram(
		histogramName,
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1),
	)
}

// DefaultSecondsHistogramView
// need register in sdkmetric.MeterProvider
// eg:
// view := SecondsHistogramView()
// mp := sdkmetric.NewMeterProvider(sdkmetric.WithView(view))
// otel.SetMeterProvider(mp)
func DefaultSecondsHistogramView(histogramName string) metricsdk.View {
	return func(instrument metricsdk.Instrument) (metricsdk.Stream, bool) {
		if instrument.Name == histogramName {
			return metricsdk.Stream{
				Name:        instrument.Name,
				Description: instrument.Description,
				Unit:        instrument.Unit,
				Aggregation: metricsdk.AggregationExplicitBucketHistogram{
					Boundaries: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1},
					NoMinMax:   true,
				},
				AttributeFilter: func(attribute.KeyValue) bool {
					return true
				},
			}, true
		}
		return metricsdk.Stream{}, false
	}
}

type options struct {
	// counter: <client/server>_requests_code_total{kind, operation, code, reason}
	requests metric.Int64Counter
	// histogram: <client/server>_requests_seconds_bucket{kind, operation}
	seconds metric.Float64Histogram
}

// Server is middleware server-side metrics.
func Server(opts ...Option) middleware.Middleware {
	op := options{}
	for _, o := range opts {
		o(&op)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// if requests and seconds are nil, return directly
			if op.requests == nil && op.seconds == nil {
				return handler(ctx, req)
			}

			var (
				code      int
				reason    string
				kind      string
				operation string
			)

			// default code
			code = status.FromGRPCCode(codes.OK)

			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = int(se.Code)
				reason = se.Reason
			}
			if op.requests != nil {
				op.requests.Add(
					ctx, 1,
					metric.WithAttributes(
						attribute.String(metricLabelKind, kind),
						attribute.String(metricLabelOperation, operation),
						attribute.Int(metricLabelCode, code),
						attribute.String(metricLabelReason, reason),
					),
				)
			}
			if op.seconds != nil {
				op.seconds.Record(
					ctx, time.Since(startTime).Seconds(),
					metric.WithAttributes(
						attribute.String(metricLabelKind, kind),
						attribute.String(metricLabelOperation, operation),
					),
				)
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
			var (
				code      int
				reason    string
				kind      string
				operation string
			)

			// default code
			code = status.FromGRPCCode(codes.OK)

			startTime := time.Now()
			if info, ok := transport.FromClientContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			reply, err := handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = int(se.Code)
				reason = se.Reason
			}
			if op.requests != nil {
				op.requests.Add(
					ctx, 1,
					metric.WithAttributes(
						attribute.String(metricLabelKind, kind),
						attribute.String(metricLabelOperation, operation),
						attribute.Int(metricLabelCode, code),
						attribute.String(metricLabelReason, reason),
					),
				)
			}
			if op.seconds != nil {
				op.seconds.Record(
					ctx, time.Since(startTime).Seconds(),
					metric.WithAttributes(
						attribute.String(metricLabelKind, kind),
						attribute.String(metricLabelOperation, operation),
					),
				)
			}
			return reply, err
		}
	}
}
