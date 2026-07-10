package log

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/trace"
)

// Option configures the OpenTelemetry slog bridge.
type Option func(*options)

type options struct {
	otel []otelslog.Option
}

// WithLoggerProvider configures the OpenTelemetry LoggerProvider.
func WithLoggerProvider(provider otellog.LoggerProvider) Option {
	return func(c *options) {
		c.otel = append(c.otel, otelslog.WithLoggerProvider(provider))
	}
}

// WithSchemaURL configures the semantic convention schema URL.
func WithSchemaURL(schemaURL string) Option {
	return func(c *options) {
		c.otel = append(c.otel, otelslog.WithSchemaURL(schemaURL))
	}
}

// WithSource configures whether source locations are emitted.
func WithSource(source bool) Option {
	return func(c *options) {
		c.otel = append(c.otel, otelslog.WithSource(source))
	}
}

// WithVersion configures the instrumentation version.
func WithVersion(version string) Option {
	return func(c *options) {
		c.otel = append(c.otel, otelslog.WithVersion(version))
	}
}

// NewHandler returns a slog handler that sends records to OpenTelemetry Logs
// and adds trace correlation attrs from the log context.
func NewHandler(name string, opts ...Option) slog.Handler {
	cfg := newOptions(opts)
	return newHandler(name, cfg)
}

func newOptions(opts []Option) options {
	var cfg options
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}

func newHandler(name string, cfg options) slog.Handler {
	return &traceHandler{next: otelslog.NewHandler(name, cfg.otel...)}
}

type traceHandler struct {
	next slog.Handler
}

func (h *traceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *traceHandler) Handle(ctx context.Context, record slog.Record) error {
	attrs := traceAttrs(ctx)
	if len(attrs) > 0 {
		record = record.Clone()
		record.AddAttrs(attrs...)
	}
	return h.next.Handle(ctx, record)
}

func (h *traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &traceHandler{next: h.next.WithAttrs(attrs)}
}

func (h *traceHandler) WithGroup(name string) slog.Handler {
	return &traceHandler{next: h.next.WithGroup(name)}
}

func traceAttrs(ctx context.Context) []slog.Attr {
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
