package log

import (
	"context"
	"fmt"
	"log"
	"os"
)

type logger struct {
	name string
	opts Options
	log  *log.Logger
}

func (l *logger) Verbose(v int) Verbose {
	switch {
	case v < 0:
		return Verbose(false)
	case verbose < v:
		return Verbose(false)
	default:
		return Verbose(true)
	}
}

func (l *logger) Print(ctx context.Context, lv Level, a ...interface{}) {
	if l.opts.Level.Enabled(lv) {
		l.log.Println(a...)
	}
}

func (l *logger) Printf(ctx context.Context, lv Level, format string, a ...interface{}) {
	if l.opts.Level.Enabled(lv) {
		l.log.Println(fmt.Sprintf(format, a...))
	}
}

func (l *logger) Printw(ctx context.Context, lv Level, kvpair ...interface{}) {
	if l.opts.Level.Enabled(lv) {
		l.log.Println(kvpair...)
	}
}

// GetLogger returns a logger instance with package name.
func GetLogger(name string, opts ...Option) Logger {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return &logger{
		name: name,
		opts: options,
		log:  log.New(os.Stdout, name, 0),
	}
}
