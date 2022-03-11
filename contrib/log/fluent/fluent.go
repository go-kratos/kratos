package fluent

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/SeeMusic/kratos/v2/log"
)

var _ log.Logger = (*Logger)(nil)

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

// Logger is fluent logger sdk.
type Logger struct {
	opts options
	log  *fluent.Fluent
}

// NewLogger new a std logger with options.
// target:
//   tcp://127.0.0.1:24224
//   unix://var/run/fluent/fluent.sock
func NewLogger(addr string, opts ...Option) (*Logger, error) {
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
	return &Logger{
		opts: option,
		log:  fl,
	}, nil
}

// Log print the kv pairs log.
func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}

	data := make(map[string]string, len(keyvals)/2+1)

	for i := 0; i < len(keyvals); i += 2 {
		data[fmt.Sprint(keyvals[i])] = fmt.Sprint(keyvals[i+1])
	}

	if err := l.log.Post(level.String(), data); err != nil {
		println(err)
	}
	return nil
}

// Close close the logger.
func (l *Logger) Close() error {
	return l.log.Close()
}
