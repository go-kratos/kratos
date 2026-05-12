package fluent

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"

	klog "github.com/go-kratos/kratos/v2/log"
)

// Option is fluentd logger option.
type Option func(*options)

type options struct {
	timeout            time.Duration
	writeTimeout       time.Duration
	bufferLimit        int
	retryWait          int
	maxRetry           int
	maxRetryWait       int
	tagPrefix          string
	async              bool
	forceStopAsyncSend bool
}

// WithTimeout with config Timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

// WithWriteTimeout with config WriteTimeout.
func WithWriteTimeout(writeTimeout time.Duration) Option {
	return func(opts *options) {
		opts.writeTimeout = writeTimeout
	}
}

// WithBufferLimit with config BufferLimit.
func WithBufferLimit(bufferLimit int) Option {
	return func(opts *options) {
		opts.bufferLimit = bufferLimit
	}
}

// WithRetryWait with config RetryWait.
func WithRetryWait(retryWait int) Option {
	return func(opts *options) {
		opts.retryWait = retryWait
	}
}

// WithMaxRetry with config MaxRetry.
func WithMaxRetry(maxRetry int) Option {
	return func(opts *options) {
		opts.maxRetry = maxRetry
	}
}

// WithMaxRetryWait with config MaxRetryWait.
func WithMaxRetryWait(maxRetryWait int) Option {
	return func(opts *options) {
		opts.maxRetryWait = maxRetryWait
	}
}

// WithTagPrefix with config TagPrefix.
func WithTagPrefix(tagPrefix string) Option {
	return func(opts *options) {
		opts.tagPrefix = tagPrefix
	}
}

// WithAsync with config Async.
func WithAsync(async bool) Option {
	return func(opts *options) {
		opts.async = async
	}
}

// WithForceStopAsyncSend with config ForceStopAsyncSend.
func WithForceStopAsyncSend(forceStopAsyncSend bool) Option {
	return func(opts *options) {
		opts.forceStopAsyncSend = forceStopAsyncSend
	}
}

// Handler writes slog records to fluentd.
type Handler struct {
	opts   options
	log    *fluent.Fluent
	attrs  []groupedAttr
	groups []string
}

type groupedAttr struct {
	groups []string
	attr   slog.Attr
}

// NewHandler returns a slog handler backed by fluentd.
// target:
//
//	tcp://127.0.0.1:24224
//	unix://var/run/fluent/fluent.sock
func NewHandler(addr string, opts ...Option) (*Handler, error) {
	option := options{}
	for _, o := range opts {
		o(&option)
	}
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	c := fluent.Config{
		Timeout:            option.timeout,
		WriteTimeout:       option.writeTimeout,
		BufferLimit:        option.bufferLimit,
		RetryWait:          option.retryWait,
		MaxRetry:           option.maxRetry,
		MaxRetryWait:       option.maxRetryWait,
		TagPrefix:          option.tagPrefix,
		Async:              option.async,
		ForceStopAsyncSend: option.forceStopAsyncSend,
	}
	switch u.Scheme {
	case "tcp":
		host, port, err2 := net.SplitHostPort(u.Host)
		if err2 != nil {
			return nil, err2
		}
		if c.FluentPort, err = strconv.Atoi(port); err != nil {
			return nil, err
		}
		c.FluentNetwork = u.Scheme
		c.FluentHost = host
	case "unix":
		c.FluentNetwork = u.Scheme
		c.FluentSocketPath = u.Path
	default:
		return nil, fmt.Errorf("unknown network: %s", u.Scheme)
	}
	fl, err := fluent.New(c)
	if err != nil {
		return nil, err
	}
	return &Handler{
		opts: option,
		log:  fl,
	}, nil
}

// NewLogger returns a slog logger backed by fluentd.
func NewLogger(addr string, opts ...Option) (*slog.Logger, error) {
	handler, err := NewHandler(addr, opts...)
	if err != nil {
		return nil, err
	}
	return klog.NewLogger(klog.WithHandler(handler)), nil
}

func (h *Handler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	data := make(map[string]string, len(h.attrs)+record.NumAttrs()+1)
	if record.Message != "" {
		data["msg"] = record.Message
	}
	for _, attr := range h.attrs {
		appendAttr(data, attr.groups, attr.attr)
	}
	record.Attrs(func(attr slog.Attr) bool {
		appendAttr(data, h.groups, attr)
		return true
	})
	if err := h.log.Post(levelString(record.Level), data); err != nil {
		println(err)
	}
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

// Close close the logger.
func (h *Handler) Close() error {
	return h.log.Close()
}

func appendAttr(data map[string]string, groups []string, attr slog.Attr) {
	attr.Value = attr.Value.Resolve()
	if attr.Value.Kind() == slog.KindGroup {
		nextGroups := groups
		if attr.Key != "" {
			nextGroups = append(append([]string{}, groups...), attr.Key)
		}
		for _, groupAttr := range attr.Value.Group() {
			appendAttr(data, nextGroups, groupAttr)
		}
		return
	}
	key := attr.Key
	if len(groups) > 0 {
		key = strings.Join(append(append([]string{}, groups...), key), ".")
	}
	data[key] = fmt.Sprint(attr.Value.Any())
}

func levelString(level slog.Level) string {
	switch {
	case level <= slog.LevelDebug:
		return "DEBUG"
	case level < slog.LevelWarn:
		return "INFO"
	case level < slog.LevelError:
		return "WARN"
	case level < slog.LevelError+4:
		return "ERROR"
	default:
		return "FATAL"
	}
}
