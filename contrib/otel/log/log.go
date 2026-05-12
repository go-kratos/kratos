package log

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/trace"

	klog "github.com/go-kratos/kratos/v2/log"
)

// Option configures the OpenTelemetry slog bridge and the Kratos log builder
// wrapping it.
type Option func(*options)

type options struct {
	otel []otelslog.Option
	log  []klog.Option
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

// WithAttrs attaches attrs to every record produced by the logger.
func WithAttrs(attrs ...slog.Attr) Option {
	return WithLogOptions(klog.WithAttrs(attrs...))
}

// WithFilter applies the provided filter options on top of the composed
// handler.
func WithFilter(opts ...klog.FilterOption) Option {
	return WithLogOptions(klog.WithFilter(opts...))
}

// WithExtractor appends attrs extracted from each log call context.
func WithExtractor(extractors ...klog.Extractor) Option {
	return WithLogOptions(klog.WithExtractor(extractors...))
}

// WithLogOptions applies Kratos core log builder options around the
// OpenTelemetry handler.
func WithLogOptions(opts ...klog.Option) Option {
	return func(c *options) {
		c.log = append(c.log, opts...)
	}
}

// NewHandler returns a slog handler that sends records to OpenTelemetry Logs.
func NewHandler(name string, opts ...Option) slog.Handler {
	cfg := newOptions(opts)
	return newHandler(name, cfg)
}

// NewLogger returns a slog logger backed by an OpenTelemetry Logs bridge handler
// and trace correlation attrs from the log context.
func NewLogger(name string, opts ...Option) *slog.Logger {
	cfg := newOptions(opts)
	logOpts := append([]klog.Option{}, cfg.log...)
	logOpts = append(logOpts,
		klog.WithExtractor(TraceAttrs),
		klog.WithHandler(otelslog.NewHandler(name, cfg.otel...)),
	)
	return klog.NewLogger(logOpts...)
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
	logOpts := append([]klog.Option{}, cfg.log...)
	logOpts = append(logOpts, klog.WithHandler(otelslog.NewHandler(name, cfg.otel...)))
	return klog.NewHandler(logOpts...)
}
