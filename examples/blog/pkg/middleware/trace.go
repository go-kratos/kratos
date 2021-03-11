package tracing

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	oteltrace "go.opentelemetry.io/otel/trace"
)

var _ propagation.TextMapCarrier = &MetadataCarrier{}

// Option is tracing option.
type Option func(*options)

type options struct {
	TracerProvider oteltrace.TracerProvider
	Propagators    propagation.TextMapPropagator
}

func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return func(opts *options) {
		opts.Propagators = propagators
	}
}

func WithTracerProvider(provider oteltrace.TracerProvider) Option {
	return func(opts *options) {
		opts.TracerProvider = provider
	}
}

type MetadataCarrier struct {
	md *metadata.MD
}

func (mc MetadataCarrier) Get(key string) string {
	values := mc.md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// Set stores the key-value pair.
func (mc MetadataCarrier) Set(key string, value string) {
	mc.md.Set(key, value)
}

// Server returns a new server middleware for OpenTelemetry.
func Server(opts ...Option) middleware.Middleware {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	if options.TracerProvider == nil {
		options.TracerProvider = otel.GetTracerProvider()
	}
	if options.TracerProvider == nil {
		options.TracerProvider = otel.GetTracerProvider()
	}
	tracer := options.TracerProvider.Tracer(
		"default",
	)
	if options.Propagators == nil {
		options.Propagators = otel.GetTextMapPropagator()
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				component string
				operation string
			)
			if info, ok := http.FromServerContext(ctx); ok {
				// HTTP span
				component = "HTTP"
				operation = info.Request.RequestURI
				ctx = propagation.NewCompositeTextMapPropagator(options.Propagators).Extract(ctx, info.Request.Header)
			} else if info, ok := grpc.FromServerContext(ctx); ok {
				// gRPC span
				component = "gRPC"
				operation = info.FullMethod
				if md, ok := metadata.FromIncomingContext(ctx); ok {
					ctx = propagation.NewCompositeTextMapPropagator(options.Propagators).Extract(ctx, MetadataCarrier{md: &md})
				}
			}
			ctx, span := tracer.Start(ctx,
				operation,
				oteltrace.WithAttributes(label.String("component", component)),
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			)
			defer span.End()
			if reply, err = handler(ctx, req); err != nil {
				span.RecordError(err)
				span.SetAttributes(
					label.String("event", "error"),
					label.String("message", err.Error()),
				)
			}
			return
		}
	}
}
