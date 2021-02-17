package tracing

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/metadata"
)

// Option is tracing option.
type Option func(*options)

type options struct {
	tracer opentracing.Tracer
}

// WithTracer sets a custom tracer to be used for this middleware, otherwise the opentracing.GlobalTracer is used.
func WithTracer(tracer opentracing.Tracer) Option {
	return func(o *options) {
		o.tracer = tracer
	}
}

// Server returns a new server middleware for OpenTracing.
func Server(opts ...Option) middleware.Middleware {
	options := options{
		tracer: opentracing.GlobalTracer(),
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				component   string
				operation   string
				spanContext opentracing.SpanContext
			)
			if info, ok := http.FromServerContext(ctx); ok {
				// HTTP span
				component = "HTTP"
				operation = info.Request.RequestURI
				spanContext, _ = options.tracer.Extract(
					opentracing.HTTPHeaders,
					opentracing.HTTPHeadersCarrier(info.Request.Header),
				)
			} else if info, ok := grpc.FromServerContext(ctx); ok {
				// gRPC span
				component = "gRPC"
				operation = info.FullMethod
				if md, ok := metadata.FromIncomingContext(ctx); ok {
					spanContext, _ = options.tracer.Extract(
						opentracing.HTTPHeaders,
						opentracing.HTTPHeadersCarrier(md),
					)
				}
			}
			span := options.tracer.StartSpan(
				operation,
				ext.RPCServerOption(spanContext),
				opentracing.Tag{Key: string(ext.Component), Value: component},
			)
			defer span.Finish()
			if reply, err = handler(ctx, req); err != nil {
				ext.Error.Set(span, true)
				span.LogFields(
					log.String("event", "error"),
					log.String("message", err.Error()),
				)
			}
			return
		}
	}
}

// Client returns a new client middleware for OpenTracing.
func Client(opts ...Option) middleware.Middleware {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				operation string
				parent    opentracing.SpanContext
				carrier   opentracing.HTTPHeadersCarrier
			)
			if span := opentracing.SpanFromContext(ctx); span != nil {
				parent = span.Context()
			}
			if info, ok := http.FromClientContext(ctx); ok {
				// HTTP span
				operation = info.Request.RequestURI
				carrier = opentracing.HTTPHeadersCarrier(info.Request.Header)
			} else if info, ok := grpc.FromClientContext(ctx); ok {
				// gRPC span
				operation = info.FullMethod
				if md, ok := metadata.FromOutgoingContext(ctx); ok {
					carrier = opentracing.HTTPHeadersCarrier(md)
					ctx = metadata.NewOutgoingContext(ctx, md)
				}
			}
			span := options.tracer.StartSpan(
				operation,
				opentracing.ChildOf(parent),
				ext.SpanKindRPCClient,
			)
			defer span.Finish()
			options.tracer.Inject(span.Context(), opentracing.HTTPHeaders, carrier)
			ctx = opentracing.ContextWithSpan(ctx, span)
			// send handler
			if reply, err = handler(ctx, req); err != nil {
				ext.Error.Set(span, true)
				span.LogFields(
					log.String("event", "error"),
					log.String("message", err.Error()),
				)
			}
			return
		}
	}
}
