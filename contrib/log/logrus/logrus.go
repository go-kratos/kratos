package logrus

import (
	"context"
	"log/slog"
	"strings"

	"github.com/sirupsen/logrus"

	klog "github.com/go-kratos/kratos/v3/log"
)

// Handler writes slog records to a logrus logger.
type Handler struct {
	logger *logrus.Logger
	attrs  logrus.Fields
	groups []string
}

// NewHandler returns a slog handler backed by logger.
func NewHandler(logger *logrus.Logger) slog.Handler {
	if logger == nil {
		logger = logrus.New()
	}
	return &Handler{logger: logger}
}

// NewLogger returns a slog logger backed by logger.
func NewLogger(logger *logrus.Logger, opts ...klog.Option) *slog.Logger {
	logOptions := append([]klog.Option{}, opts...)
	logOptions = append(logOptions, klog.WithHandler(NewHandler(logger)))
	return klog.NewLogger(logOptions...)
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return toLogrusLevel(level) <= h.logger.Level
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	fields := make(logrus.Fields, len(h.attrs)+record.NumAttrs())
	for key, value := range h.attrs {
		fields[key] = value
	}
	record.Attrs(func(attr slog.Attr) bool {
		appendAttr(fields, h.groups, attr)
		return true
	})
	entry := h.logger.WithFields(fields)
	entry.Log(toLogrusLevel(record.Level), record.Message)
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := *h
	next.attrs = make(logrus.Fields, len(h.attrs)+len(attrs))
	for key, value := range h.attrs {
		next.attrs[key] = value
	}
	for _, attr := range attrs {
		appendAttr(next.attrs, h.groups, attr)
	}
	return &next
}

func (h *Handler) WithGroup(name string) slog.Handler {
	next := *h
	next.groups = append(append([]string{}, h.groups...), name)
	return &next
}

func appendAttr(fields logrus.Fields, groups []string, attr slog.Attr) {
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
	fields[key] = attr.Value.Any()
}

func toLogrusLevel(level slog.Level) logrus.Level {
	switch {
	case level <= slog.LevelDebug:
		return logrus.DebugLevel
	case level < slog.LevelWarn:
		return logrus.InfoLevel
	case level < slog.LevelError:
		return logrus.WarnLevel
	case level < slog.LevelError+4:
		return logrus.ErrorLevel
	default:
		return logrus.FatalLevel
	}
}
