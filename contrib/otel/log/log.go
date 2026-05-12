package log

import (
	"context"
	"log/slog"

	klog "github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/trace"
)

// Option configures the OpenTelemetry slog bridge handler. Kratos logger
// options are applied by composing [NewHandler] with the core log builder.
type Option = otelslog.Option

// WithLoggerProvider configures the OpenTelemetry LoggerProvider.
var WithLoggerProvider = otelslog.WithLoggerProvider

// WithAttributes configures the instrumentation scope attributes.
var WithAttributes = otelslog.WithAttributes

// WithSchemaURL configures the semantic convention schema URL.
var WithSchemaURL = otelslog.WithSchemaURL

// WithSource configures whether source locations are emitted.
var WithSource = otelslog.WithSource

// WithVersion configures the instrumentation version.
var WithVersion = otelslog.WithVersion

// NewHandler returns a slog handler that sends records to OpenTelemetry Logs.
func NewHandler(name string, opts ...Option) slog.Handler {
	return otelslog.NewHandler(name, opts...)
}

// NewLogger returns a slog logger backed by an OpenTelemetry Logs bridge handler
// and trace correlation attrs from the log context.
func NewLogger(name string, opts ...Option) *slog.Logger {
	return klog.NewLogger(
		klog.WithHandler(NewHandler(name, opts...)),
		klog.WithExtractor(TraceAttrs),
	)
}

// TraceAttrs pulls trace_id / span_id / trace_flags from ctx.
func TraceAttrs(ctx context.Context) []slog.Attr {
	span := trace.SpanContextFromContext(ctx)
	if !span.IsValid() {
		return nil
	}
	return []slog.Attr{
		slog.String("trace_id", span.TraceID().String()),
		slog.String("span_id", span.SpanID().String()),
		slog.String("trace_flags", span.TraceFlags().String()),
	}
}
