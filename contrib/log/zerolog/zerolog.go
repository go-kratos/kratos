package zerolog

import (
	"context"
	"log/slog"
	"strings"

	"github.com/rs/zerolog"

	klog "github.com/go-kratos/kratos/v3/log"
)

// Handler writes slog records to a zerolog logger.
type Handler struct {
	logger *zerolog.Logger
	attrs  []groupedAttr
	groups []string
}

type groupedAttr struct {
	groups []string
	attr   slog.Attr
}

// NewHandler returns a slog handler backed by logger.
func NewHandler(logger *zerolog.Logger) slog.Handler {
	if logger == nil {
		l := zerolog.Nop()
		logger = &l
	}
	return &Handler{logger: logger}
}

// NewLogger returns a slog logger backed by logger.
func NewLogger(logger *zerolog.Logger, opts ...klog.Option) *slog.Logger {
	logOptions := append([]klog.Option{}, opts...)
	logOptions = append(logOptions, klog.WithHandler(NewHandler(logger)))
	return klog.NewLogger(logOptions...)
}

func (h *Handler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	event := h.event(record.Level)
	if event == nil {
		return nil
	}
	for _, attr := range h.attrs {
		appendAttr(event, attr.groups, attr.attr)
	}
	record.Attrs(func(attr slog.Attr) bool {
		appendAttr(event, h.groups, attr)
		return true
	})
	event.Msg(record.Message)
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	next := *h
	next.attrs = append([]groupedAttr{}, h.attrs...)
	for _, attr := range attrs {
		next.attrs = append(next.attrs, groupedAttr{
			groups: append([]string{}, h.groups...),
			attr:   attr,
		})
	}
	return &next
}

func (h *Handler) WithGroup(name string) slog.Handler {
	next := *h
	next.groups = append(append([]string{}, h.groups...), name)
	return &next
}

func (h *Handler) event(level slog.Level) *zerolog.Event {
	switch {
	case level <= slog.LevelDebug:
		return h.logger.Debug()
	case level < slog.LevelWarn:
		return h.logger.Info()
	case level < slog.LevelError:
		return h.logger.Warn()
	case level < slog.LevelError+4:
		return h.logger.Error()
	default:
		return h.logger.Fatal()
	}
}

func appendAttr(event *zerolog.Event, groups []string, attr slog.Attr) {
	attr.Value = attr.Value.Resolve()
	if attr.Value.Kind() == slog.KindGroup {
		nextGroups := groups
		if attr.Key != "" {
			nextGroups = append(append([]string{}, groups...), attr.Key)
		}
		for _, groupAttr := range attr.Value.Group() {
			appendAttr(event, nextGroups, groupAttr)
		}
		return
	}
	key := attr.Key
	if len(groups) > 0 {
		key = strings.Join(append(append([]string{}, groups...), key), ".")
	}
	event.Any(key, attr.Value.Any())
}
