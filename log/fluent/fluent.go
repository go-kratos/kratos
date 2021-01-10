package fluent

import (
	"fmt"
	"net"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/go-kratos/kratos/v2/log"
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
	skip               int
}

// Timeout with config Timeout.
func Timeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

// WriteTimeout with config WriteTimeout.
func WriteTimeout(writeTimeout time.Duration) Option {
	return func(opts *options) {
		opts.writeTimeout = writeTimeout
	}
}

// BufferLimit with config BufferLimit.
func BufferLimit(bufferLimit int) Option {
	return func(opts *options) {
		opts.bufferLimit = bufferLimit
	}
}

// RetryWait with config RetryWait.
func RetryWait(retryWait int) Option {
	return func(opts *options) {
		opts.retryWait = retryWait
	}
}

// MaxRetry with config MaxRetry.
func MaxRetry(maxRetry int) Option {
	return func(opts *options) {
		opts.maxRetry = maxRetry
	}
}

// MaxRetryWait with config MaxRetryWait.
func MaxRetryWait(maxRetryWait int) Option {
	return func(opts *options) {
		opts.maxRetryWait = maxRetryWait
	}
}

// TagPrefix with config TagPrefix.
func TagPrefix(tagPrefix string) Option {
	return func(opts *options) {
		opts.tagPrefix = tagPrefix
	}
}

// Async with config Async.
func Async(async bool) Option {
	return func(opts *options) {
		opts.async = async
	}
}

// ForceStopAsyncSend with config ForceStopAsyncSend.
func ForceStopAsyncSend(forceStopAsyncSend bool) Option {
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
//   unix:///var/run/fluent/fluent.sock
func NewLogger(addr string, opts ...Option) (*Logger, error) {
	options := options{skip: 4}
	for _, o := range opts {
		o(&options)
	}
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	c := fluent.Config{
		Timeout:            options.timeout,
		WriteTimeout:       options.writeTimeout,
		BufferLimit:        options.bufferLimit,
		RetryWait:          options.retryWait,
		MaxRetry:           options.maxRetry,
		MaxRetryWait:       options.maxRetryWait,
		TagPrefix:          options.tagPrefix,
		Async:              options.async,
		ForceStopAsyncSend: options.forceStopAsyncSend,
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
	return &Logger{
		opts: options,
		log:  fl,
	}, nil
}

func (l *Logger) stackTrace(path string) string {
	idx := strings.LastIndexByte(path, '/')
	if idx == -1 {
		return path
	}
	idx = strings.LastIndexByte(path[:idx], '/')
	if idx == -1 {
		return path
	}
	return path[idx+1:]
}

// Print print the kv pairs log.
func (l *Logger) Print(kvpair ...interface{}) {
	if len(kvpair) == 0 {
		return
	}
	if len(kvpair)%2 != 0 {
		kvpair = append(kvpair, "")
	}

	data := make(map[string]string, len(kvpair)/2+1)
	if _, file, line, ok := runtime.Caller(l.opts.skip); ok {
		data[l.stackTrace(file)] = strconv.Itoa(line)
	}
	for i := 0; i < len(kvpair); i += 2 {
		data[fmt.Sprint(kvpair[i])] = fmt.Sprint(kvpair[i+1])
	}

	if err := l.log.Post(data["module"], data); err != nil {
		println(err)
	}
}

// Close close the logger.
func (l *Logger) Close() error {
	return l.log.Close()
}
