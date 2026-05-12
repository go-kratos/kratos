package zap

import (
	"context"
	"log/slog"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	klog "github.com/go-kratos/kratos/v3/log"
)

// Handler writes slog records to a zap logger.
type Handler struct {
	logger *zap.Logger
	attrs  []zap.Field
	groups []string
}

// NewHandler returns a slog handler backed by zlog.
func NewHandler(zlog *zap.Logger) slog.Handler {
	if zlog == nil {
		zlog = zap.NewNop()
	}
	return &Handler{logger: zlog}
}

// NewLogger returns a slog logger backed by zlog.
func NewLogger(zlog *zap.Logger, opts ...klog.Option) *slog.Logger {
	logOptions := append([]klog.Option{}, opts...)
	logOptions = append(logOptions, klog.WithHandler(NewHandler(zlog)))
	return klog.NewLogger(logOptions...)
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	zapLevel := toZapLevel(level)
	return zapLevel >= zapcore.DPanicLevel || h.logger.Core().Enabled(zapLevel)
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	zapLevel := toZapLevel(record.Level)
	checked := h.logger.Check(zapLevel, record.Message)
	if checked == nil {
		return nil
	}
	fields := make([]zap.Field, 0, len(h.attrs)+record.NumAttrs())
	fields = append(fields, h.attrs...)
	record.Attrs(func(attr slog.Attr) bool {
		appendAttr(&fields, h.groups, attr)
		return true
	})
	checked.Write(fields...)
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := *h
	next.attrs = append(append([]zap.Field{}, h.attrs...), attrsToFields(h.groups, attrs)...)
	return &next
}

func (h *Handler) WithGroup(name string) slog.Handler {
	next := *h
	next.groups = append(append([]string{}, h.groups...), name)
	return &next
}

func attrsToFields(groups []string, attrs []slog.Attr) []zap.Field {
	fields := make([]zap.Field, 0, len(attrs))
	for _, attr := range attrs {
		appendAttr(&fields, groups, attr)
	}
	return fields
}

func appendAttr(fields *[]zap.Field, groups []string, attr slog.Attr) {
	attr.Value = attr.Value.Resolve()
	if attr.Value.Kind() == slog.KindGroup {
		nextGroups := groups
		if attr.Key != "" {
			nextGroups = append(append([]string{}, groups...), attr.Key)
		}
		for _, groupAttr := range attr.Value.Group() {
			appendAttr(fields, nextGroups, groupAttr)
		}
		return
	}
	key := attr.Key
	if len(groups) > 0 {
		key = strings.Join(append(append([]string{}, groups...), key), ".")
	}
	*fields = append(*fields, zap.Any(key, attr.Value.Any()))
}

func toZapLevel(level slog.Level) zapcore.Level {
	switch {
	case level <= slog.LevelDebug:
		return zapcore.DebugLevel
	case level < slog.LevelWarn:
		return zapcore.InfoLevel
	case level < slog.LevelError:
		return zapcore.WarnLevel
	case level < slog.LevelError+4:
		return zapcore.ErrorLevel
	default:
		return zapcore.FatalLevel
	}
}

// Sync flushes buffered log entries.
func (h *Handler) Sync() error {
	return h.logger.Sync()
}

// Close flushes buffered log entries.
func (h *Handler) Close() error {
	return h.Sync()
}
