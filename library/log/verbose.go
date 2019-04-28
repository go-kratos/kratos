package log

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// V reports whether verbosity at the call site is at least the requested level.
// The returned value is a boolean of type Verbose, which implements Info, Infov etc.
// These methods will write to the Info log if called.
// Thus, one may write either
//	if log.V(2) { log.Info("log this") }
// or
//	log.V(2).Info("log this")
// The second form is shorter but the first is cheaper if logging is off because it does
// not evaluate its arguments.
//
// Whether an individual call to V generates a log record depends on the setting of
// the Config.VLevel and Config.Module flags; both are off by default. If the level in the call to
// V is at least the value of Config.VLevel, or of Config.Module for the source file containing the
// call, the V call will log.
// v must be more than 0.
func V(v int32) Verbose {
	var (
		file string
	)
	if v < 0 {
		return Verbose(false)
	} else if c.V >= v {
		return Verbose(true)
	}
	if pc, _, _, ok := runtime.Caller(1); ok {
		file, _ = runtime.FuncForPC(pc).FileLine(pc)
	}
	if strings.HasSuffix(file, ".go") {
		file = file[:len(file)-3]
	}
	if slash := strings.LastIndex(file, "/"); slash >= 0 {
		file = file[slash+1:]
	}
	for filter, lvl := range c.Module {
		var match bool
		if match = filter == file; !match {
			match, _ = filepath.Match(filter, file)
		}
		if match {
			return Verbose(lvl >= v)
		}
	}
	return Verbose(false)
}

// Info logs a message at the info log level.
func (v Verbose) Info(format string, args ...interface{}) {
	if v {
		h.Log(context.Background(), _infoLevel, KV(_log, fmt.Sprintf(format, args...)))
	}
}

// Infov logs a message at the info log level.
func (v Verbose) Infov(ctx context.Context, args ...D) {
	if v {
		h.Log(ctx, _infoLevel, args...)
	}
}

// Infow logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func (v Verbose) Infow(ctx context.Context, args ...interface{}) {
	if v {
		h.Log(ctx, _infoLevel, logw(args)...)
	}
}

// Close close resource.
func (v Verbose) Close() (err error) {
	if h == nil {
		return
	}
	return h.Close()
}
