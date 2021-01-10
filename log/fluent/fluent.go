package fluent

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/go-kratos/kratos/v2/log"
)

var _ log.Logger = (*fluentLogger)(nil)

// Option is fluentd logger option.
type Option func(*options)

type options struct {
	Timeout            time.Duration
	WriteTimeout       time.Duration
	BufferLimit        int
	RetryWait          int
	MaxRetry           int
	MaxRetryWait       int
	TagPrefix          string
	Async              bool
	ForceStopAsyncSend bool
}

// Timeout with config Timeout.
func Timeout(val time.Duration) Option {
	return func(opts *options) {
		opts.Timeout = val
	}
}

// WriteTimeout with config WriteTimeout.
func WriteTimeout(val time.Duration) Option {
	return func(opts *options) {
		opts.WriteTimeout = val
	}
}

// BufferLimit with config BufferLimit.
func BufferLimit(val int) Option {
	return func(opts *options) {
		opts.BufferLimit = val
	}
}

// RetryWait with config RetryWait.
func RetryWait(val int) Option {
	return func(opts *options) {
		opts.RetryWait = val
	}
}

// MaxRetry with config MaxRetry.
func MaxRetry(val int) Option {
	return func(opts *options) {
		opts.MaxRetry = val
	}
}

// MaxRetryWait with config MaxRetryWait.
func MaxRetryWait(val int) Option {
	return func(opts *options) {
		opts.MaxRetryWait = val
	}
}

// TagPrefix with config TagPrefix.
func TagPrefix(val string) Option {
	return func(opts *options) {
		opts.TagPrefix = val
	}
}

// Async with config Async.
func Async(val bool) Option {
	return func(opts *options) {
		opts.Async = val
	}
}

// ForceStopAsyncSend with config ForceStopAsyncSend.
func ForceStopAsyncSend(val bool) Option {
	return func(opts *options) {
		opts.ForceStopAsyncSend = val
	}
}

type fluentLogger struct {
	opts options
	log  *fluent.Fluent
}

// NewLogger new a std logger with options.
// target:
//   tcp://127.0.0.1:24224
//   unix:///var/run/fluent/fluent.sock
func NewLogger(target string, opts ...Option) (log.Logger, error) {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	c := fluent.Config{
		Timeout:            options.Timeout,
		WriteTimeout:       options.WriteTimeout,
		BufferLimit:        options.BufferLimit,
		RetryWait:          options.RetryWait,
		MaxRetry:           options.MaxRetry,
		MaxRetryWait:       options.MaxRetryWait,
		TagPrefix:          options.TagPrefix,
		Async:              options.Async,
		ForceStopAsyncSend: options.ForceStopAsyncSend,
	}
	switch u.Scheme {
	case "tcp":
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return nil, err
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
	return &fluentLogger{
		opts: options,
		log:  fl,
	}, nil
}

func (f *fluentLogger) Print(level log.Level, kvpair ...interface{}) {
	if len(kvpair) == 0 {
		return
	}
	if len(kvpair)%2 != 0 {
		kvpair = append(kvpair, "")
	}

	data := make(map[string]string, len(kvpair)/2+1)
	data["level"] = level.String()
	for i := 0; i < len(kvpair); i += 2 {
		data[fmt.Sprint(kvpair[i])] = fmt.Sprint(kvpair[i+1])
	}

	if err := f.log.Post(data["module"], data); err != nil {
		println(err)
	}
}

func (f *fluentLogger) Close() error {
	return f.log.Close()
}
