package tracing

import (
    "context"
    "github.com/go-kratos/kratos/v2/log"
    "go.opentelemetry.io/otel/trace"
)

// TraceID returns a traceid valuer.
func TraceID() log.Valuer {
    return func(ctx context.Context) interface{} {
        if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
            return span.TraceID().String()
        }
        return ""
    }
}

// SpanID returns a spanid valuer.
func SpanID() log.Valuer {
    return func(ctx context.Context) interface{} {
        if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
            return span.SpanID().String()
        }
        return ""
    }
}
