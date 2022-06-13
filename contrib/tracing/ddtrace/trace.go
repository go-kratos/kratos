package ddtrace

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Start(env, service, version string) {
	hostname, _ := os.Hostname()
	tracer.Start(
		tracer.WithTraceEnabled(true),
		tracer.WithAgentAddr(os.Getenv("DD_AGENT_HOST")),
		tracer.WithEnv(env),
		tracer.WithHostname(hostname),
		tracer.WithService(service),
		tracer.WithServiceVersion(version),
	)
}

// Server is a server ddtracer middleware.
func Server(serviceName string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				operation string
				message   string
				kind      string
				opts      []ddtrace.StartSpanOption

				textHeader = make(map[string]string)
			)

			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
				trHeader := info.RequestHeader()
				for _, key := range trHeader.Keys() {
					textHeader[key] = trHeader.Get(key)
				}
				if parentSpanCtx, err := tracer.Extract(tracer.TextMapCarrier(textHeader)); err == nil {
					opts = append(opts, tracer.ChildOf(parentSpanCtx))
				}
			}

			opts = append(opts, []ddtrace.StartSpanOption{
				tracer.ServiceName(serviceName),
				tracer.ResourceName(operation),
				tracer.SpanType(kind),
				tracer.Measured(),
			}...)

			span, ctx := tracer.StartSpanFromContext(ctx, operation, opts...)
			defer span.Finish()

			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
				message = se.Message
			}

			if err == nil {
				span.SetTag(ext.HTTPCode, http.StatusOK)
			} else {
				span.SetTag(ext.HTTPCode, strconv.Itoa(int(code)))
				span.SetTag(ext.ErrorMsg, message)
				span.SetTag(ext.ErrorDetails, reason)
			}

			return
		}
	}
}

// Client is a client ddtracer middleware
func Client() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			span, ok := tracer.SpanFromContext(ctx)
			if ok {
				spanCtx := span.Context()
				md := metadata.New(map[string]string{
					tracer.DefaultTraceIDHeader:  strconv.FormatUint(spanCtx.TraceID(), 10),
					tracer.DefaultParentIDHeader: strconv.FormatUint(spanCtx.SpanID(), 10),
				})
				ctx = metadata.MergeToClientContext(ctx, md)
			}
			return handler(ctx, req)
		}
	}
}
