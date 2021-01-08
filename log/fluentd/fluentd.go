package fluentd

import (
	"fmt"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/go-kratos/kratos/v2/log"
)

var _ log.Logger = (*fluentdLogger)(nil)

// Option is fluentd logger option.
type Option func(*options)

type options struct {
	FluentPort         int
	FluentHost         string
	FluentNetwork      string
	FluentSocketPath   string
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

type fluentdLogger struct {
	opts options
	log  *fluent.Fluent
}

// NewLogger new a std logger with options.
func NewLogger(opts ...Option) log.Logger {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	fl, err := fluent.New(fluent.Config{
		FluentPort:         options.FluentPort,
		FluentHost:         options.FluentHost,
		FluentNetwork:      options.FluentNetwork,
		FluentSocketPath:   options.FluentSocketPath,
		Timeout:            options.Timeout,
		WriteTimeout:       options.WriteTimeout,
		BufferLimit:        options.BufferLimit,
		RetryWait:          options.RetryWait,
		MaxRetry:           options.MaxRetry,
		MaxRetryWait:       options.MaxRetryWait,
		TagPrefix:          options.TagPrefix,
		Async:              options.Async,
		ForceStopAsyncSend: options.ForceStopAsyncSend,
	})
	if err != nil {
		panic(err)
	}
	return &fluentdLogger{
		opts: options,
		log:  fl,
	}
}

func (f *fluentdLogger) Print(kvpair ...interface{}) {
	if len(kvpair) == 0 {
		return
	}
	if len(kvpair)%2 != 0 {
		kvpair = append(kvpair, "")
	}

	tag := "" // fixme: TBD
	data := make(map[string]string)
	for i := 0; i < len(kvpair); i += 2 {
		data[fmt.Sprintf("%s", kvpair[i])] = fmt.Sprintf("%s", kvpair[i+1])
	}

	err := f.log.Post(tag, data)
	if err != nil {
		println(err)
	}
}

func (f *fluentdLogger) Close() error {
	return f.log.Close()
}
