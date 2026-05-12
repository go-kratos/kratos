package log

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// Format selects the encoding used by the default handler builder.
type Format int

const (
	// FormatText writes records using [slog.NewTextHandler].
	FormatText Format = iota
	// FormatJSON writes records using [slog.NewJSONHandler].
	FormatJSON
)

// Option configures [NewHandler] / [NewLogger].
type Option func(*handlerConfig)

// Extractor extracts attrs from a log call context.
type Extractor func(context.Context) []slog.Attr

type handlerConfig struct {
	handler     slog.Handler
	writer      io.Writer
	format      Format
	level       Leveler
	addSource   bool
	replaceAttr func(groups []string, a slog.Attr) slog.Attr
	attrs       []slog.Attr
	extractors  []Extractor
	filter      []FilterOption
}

// WithHandler sets the base slog handler. When set, writer/format/level/source
// options only apply to the default handler and are ignored.
func WithHandler(h slog.Handler) Option {
	return func(c *handlerConfig) { c.handler = h }
}

// WithAttrs attaches attrs to every record produced by the logger.
func WithAttrs(attrs ...slog.Attr) Option {
	return func(c *handlerConfig) { c.attrs = append(c.attrs, attrs...) }
}

// WithExtractor appends attrs extracted from each log call context.
func WithExtractor(extractors ...Extractor) Option {
	return func(c *handlerConfig) {
		for _, e := range extractors {
			if e != nil {
				c.extractors = append(c.extractors, e)
			}
		}
	}
}

// WithWriter sets the destination writer for the base handler. Defaults to
// [os.Stderr].
func WithWriter(w io.Writer) Option {
	return func(c *handlerConfig) { c.writer = w }
}

// WithFormat selects between text and JSON encoding. Defaults to [FormatText].
func WithFormat(f Format) Option {
	return func(c *handlerConfig) { c.format = f }
}

// WithLevel sets the minimum level for the base handler.
func WithLevel(l Leveler) Option {
	return func(c *handlerConfig) { c.level = l }
}

// WithAddSource toggles inclusion of the source file/line.
func WithAddSource(b bool) Option {
	return func(c *handlerConfig) { c.addSource = b }
}

// WithReplaceAttr installs a custom ReplaceAttr on the base handler.
func WithReplaceAttr(fn func(groups []string, a slog.Attr) slog.Attr) Option {
	return func(c *handlerConfig) { c.replaceAttr = fn }
}

// WithFilter applies the provided filter options on top of the composed
// handler.
func WithFilter(opts ...FilterOption) Option {
	return func(c *handlerConfig) { c.filter = append(c.filter, opts...) }
}

// NewHandler builds a composed [slog.Handler] with kratos defaults:
//   - text encoding to stderr at LevelInfo
//   - context attrs from [ContextWithAttrs] are merged in
//
// Additional decorators are layered as configured.
func NewHandler(opts ...Option) slog.Handler {
	cfg := &handlerConfig{
		writer:     os.Stderr,
		format:     FormatText,
		level:      LevelInfo,
		extractors: []Extractor{AttrsFromContext},
	}
	for _, o := range opts {
		o(cfg)
	}
	h := cfg.handler
	if h == nil {
		h = newBaseHandler(cfg)
	}
	if len(cfg.filter) > 0 {
		h = newFilterHandler(h, cfg.filter...)
	}
	if len(cfg.attrs) > 0 {
		h = h.WithAttrs(cfg.attrs)
	}
	return newContextHandler(h, cfg.extractors...)
}

// NewLogger is shorthand for slog.New(NewHandler(opts...)).
func NewLogger(opts ...Option) *slog.Logger {
	return slog.New(NewHandler(opts...))
}

func newBaseHandler(cfg *handlerConfig) slog.Handler {
	hopts := &slog.HandlerOptions{
		Level:       cfg.level,
		AddSource:   cfg.addSource,
		ReplaceAttr: cfg.replaceAttr,
	}
	switch cfg.format {
	case FormatJSON:
		return slog.NewJSONHandler(cfg.writer, hopts)
	default:
		return slog.NewTextHandler(cfg.writer, hopts)
	}
}
