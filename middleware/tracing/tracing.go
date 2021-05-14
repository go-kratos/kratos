package tracing

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// Option is tracing option.
type Option func(*options)

type options struct {
	TracerProvider trace.TracerProvider
	Propagators    propagation.TextMapPropagator
}

func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return func(opts *options) {
		opts.Propagators = propagators
	}
}

func WithTracerProvider(provider trace.TracerProvider) Option {
	return func(opts *options) {
		opts.TracerProvider = provider
	}
}

type MetadataCarrier struct {
	md *metadata.MD
}

var _ propagation.TextMapCarrier = &MetadataCarrier{}

func (mc MetadataCarrier) Get(key string) string {
	values := mc.md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (mc MetadataCarrier) Set(key string, value string) {
	mc.md.Set(key, value)
}

func (mc MetadataCarrier) Keys() []string {
	keys := make([]string, 0, mc.md.Len())
	for key := range *mc.md {
		keys = append(keys, key)
	}
	return keys
}

// Server returns a new server middleware for OpenTelemetry.
func Server(opts ...Option) middleware.Middleware {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	if options.TracerProvider != nil {
		otel.SetTracerProvider(options.TracerProvider)
	}
	if options.Propagators != nil {
		otel.SetTextMapPropagator(options.Propagators)
	}
	tracer := otel.Tracer("server")
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				component string
				operation string
			)
			if tr, ok := transport.FromContext(ctx); ok {
				component = string(tr.Kind)
				operation = tr.Request.FullPath
				ctx = otel.GetTextMapPropagator().Extract(ctx, tr.Request.Metadata)
			}
			ctx, span := tracer.Start(ctx,
				operation,
				trace.WithAttributes(attribute.String("component", component)),
				trace.WithSpanKind(trace.SpanKindServer),
			)
			defer span.End()
			if reply, err = handler(ctx, req); err != nil {
				span.RecordError(err)
				span.SetAttributes(
					attribute.String("event", "error"),
					attribute.String("message", err.Error()),
				)
			}
			return
		}
	}
}

// Client returns a new client middleware for OpenTelemetry.
func Client(opts ...Option) middleware.Middleware {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	if options.TracerProvider != nil {
		otel.SetTracerProvider(options.TracerProvider)
	}
	if options.Propagators != nil {
		otel.SetTextMapPropagator(options.Propagators)
	}
	tracer := otel.Tracer("client")
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				component string
				operation string
				carrier   propagation.TextMapCarrier
			)
			if tr, ok := transport.FromContext(ctx); ok {
				component = string(tr.Kind)
				operation = tr.Request.FullPath
				carrier = tr.Request.Metadata
			}

			ctx, span := tracer.Start(ctx,
				operation,
				trace.WithAttributes(attribute.String("component", component)),
				trace.WithSpanKind(trace.SpanKindClient),
			)
			defer span.End()
			otel.GetTextMapPropagator().Inject(ctx, carrier)
			if reply, err = handler(ctx, req); err != nil {
				span.RecordError(err)
				span.SetAttributes(
					attribute.String("event", "error"),
					attribute.String("message", err.Error()),
				)
			}
			return
		}
	}
}
