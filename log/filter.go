package log

import (
	"context"
	"log/slog"
	"strings"
)

const redactedValue = "***"

// FilterOption configures filtering in [WithFilter].
type FilterOption func(*filterConfig)

type filterConfig struct {
	keys   map[string]struct{}
	filter func(ctx context.Context, record slog.Record) bool
}

// FilterKey redacts the values of attributes whose key matches any of the
// provided keys. Keys may be leaf names ("password") or dotted group paths
// ("user.password").
func FilterKey(keys ...string) FilterOption {
	return func(c *filterConfig) {
		if c.keys == nil {
			c.keys = make(map[string]struct{}, len(keys))
		}
		for _, k := range keys {
			c.keys[k] = struct{}{}
		}
	}
}

// FilterFunc drops records for which fn returns true. fn is evaluated after key
// redaction.
func FilterFunc(fn func(ctx context.Context, record slog.Record) bool) FilterOption {
	return func(c *filterConfig) { c.filter = fn }
}

func newFilterHandler(next slog.Handler, opts ...FilterOption) slog.Handler {
	if next == nil {
		next = discardHandler{}
	}
	cfg := &filterConfig{}
	for _, o := range opts {
		o(cfg)
	}
	if len(cfg.keys) == 0 && cfg.filter == nil {
		return next
	}
	return &filterHandler{next: next, cfg: cfg}
}

type filterHandler struct {
	next   slog.Handler
	cfg    *filterConfig
	groups []string
}

func (h *filterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *filterHandler) Handle(ctx context.Context, record slog.Record) error {
	if h.needsRewrite() {
		record = h.rewrite(record)
	}
	if h.cfg.filter != nil && h.cfg.filter(ctx, record) {
		return nil
	}
	return h.next.Handle(ctx, record)
}

func (h *filterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if h.needsRewrite() {
		attrs = h.redactAttrs(h.groups, attrs)
	}
	next := *h
	next.next = h.next.WithAttrs(attrs)
	return &next
}

func (h *filterHandler) WithGroup(name string) slog.Handler {
	next := *h
	next.groups = append(append([]string{}, h.groups...), name)
	next.next = h.next.WithGroup(name)
	return &next
}

func (h *filterHandler) needsRewrite() bool {
	return len(h.cfg.keys) > 0
}

func (h *filterHandler) rewrite(record slog.Record) slog.Record {
	cloned := record.Clone()
	attrs := make([]slog.Attr, 0, cloned.NumAttrs())
	cloned.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	redacted := h.redactAttrs(h.groups, attrs)
	out := slog.NewRecord(cloned.Time, cloned.Level, cloned.Message, cloned.PC)
	out.AddAttrs(redacted...)
	return out
}

func (h *filterHandler) redactAttrs(groups []string, attrs []slog.Attr) []slog.Attr {
	out := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		out[i] = h.redactAttr(groups, a)
	}
	return out
}

func (h *filterHandler) redactAttr(groups []string, a slog.Attr) slog.Attr {
	a.Value = a.Value.Resolve()
	if a.Value.Kind() == slog.KindGroup {
		group := a.Value.Group()
		next := make([]slog.Attr, len(group))
		nextGroups := appendPath(groups, a.Key)
		for i, ga := range group {
			next[i] = h.redactAttr(nextGroups, ga)
		}
		return slog.Attr{Key: a.Key, Value: slog.GroupValue(next...)}
	}
	if h.matchesKey(groups, a.Key) {
		return slog.Attr{Key: a.Key, Value: slog.StringValue(redactedValue)}
	}
	return a
}

func (h *filterHandler) matchesKey(groups []string, key string) bool {
	if _, ok := h.cfg.keys[key]; ok {
		return true
	}
	if len(groups) == 0 {
		return false
	}
	path := strings.Join(appendPath(groups, key), ".")
	_, ok := h.cfg.keys[path]
	return ok
}

func appendPath(groups []string, key string) []string {
	if key == "" {
		return groups
	}
	next := make([]string, 0, len(groups)+1)
	next = append(next, groups...)
	next = append(next, key)
	return next
}
