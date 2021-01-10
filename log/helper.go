package log

import (
	"fmt"
)

var nop Logger = new(nopLogger)

// Option .
type Option func(*options)

type options struct {
	level   Level
	verbose Verbose
}

// AllowLevel .
func AllowLevel(l Level) Option {
	return func(o *options) {
		o.level = l
	}
}

// AllowVerbose .
func AllowVerbose(v Verbose) Option {
	return func(o *options) {
		o.verbose = v
	}
}

// Helper is a logger helper.
type Helper struct {
	Logger

	opts options
}

// NewHelper new a logger helper.
func NewHelper(name string, l Logger, opts ...Option) *Helper {
	options := options{}
	for _, o := range opts {
		o(&options)
	}
	return &Helper{Logger: With(l, "module", name), opts: options}
}

// V logs a message at verbose level.
func (h *Helper) V(v Verbose) Logger {
	if h.opts.verbose.Enabled(v) {
		return nop
	}
	return With(h, VerboseKey, v)
}

// Debug logs a message at debug level.
func (h *Helper) Debug(a ...interface{}) {
	if h.opts.level.Enabled(LevelDebug) {
		Debug(h).Print("log", fmt.Sprint(a...))
	}
}

// Debugf logs a message at debug level.
func (h *Helper) Debugf(format string, a ...interface{}) {
	if h.opts.level.Enabled(LevelDebug) {
		Debug(h).Print("log", fmt.Sprintf(format, a...))
	}
}

// Debugw logs a message at debug level.
func (h *Helper) Debugw(kvpair ...interface{}) {
	if h.opts.level.Enabled(LevelDebug) {
		Debug(h).Print(kvpair...)
	}
}

// Info logs a message at info level.
func (h *Helper) Info(a ...interface{}) {
	if h.opts.level.Enabled(LevelInfo) {
		Info(h).Print("log", fmt.Sprint(a...))
	}
}

// Infof logs a message at info level.
func (h *Helper) Infof(format string, a ...interface{}) {
	if h.opts.level.Enabled(LevelInfo) {
		Info(h).Print("log", fmt.Sprintf(format, a...))
	}
}

// Infow logs a message at info level.
func (h *Helper) Infow(kvpair ...interface{}) {
	if h.opts.level.Enabled(LevelInfo) {
		Info(h).Print(kvpair...)
	}
}

// Warn logs a message at warn level.
func (h *Helper) Warn(a ...interface{}) {
	if h.opts.level.Enabled(LevelWarn) {
		Warn(h).Print("log", fmt.Sprint(a...))
	}
}

// Warnf logs a message at warnf level.
func (h *Helper) Warnf(format string, a ...interface{}) {
	if h.opts.level.Enabled(LevelWarn) {
		Warn(h).Print("log", fmt.Sprintf(format, a...))
	}
}

// Warnw logs a message at warnf level.
func (h *Helper) Warnw(kvpair ...interface{}) {
	if h.opts.level.Enabled(LevelWarn) {
		Warn(h).Print(kvpair...)
	}
}

// Error logs a message at error level.
func (h *Helper) Error(a ...interface{}) {
	if h.opts.level.Enabled(LevelError) {
		Error(h).Print("log", fmt.Sprint(a...))
	}
}

// Errorf logs a message at error level.
func (h *Helper) Errorf(format string, a ...interface{}) {
	if h.opts.level.Enabled(LevelError) {
		Error(h).Print("log", fmt.Sprintf(format, a...))
	}
}

// Errorw logs a message at error level.
func (h *Helper) Errorw(kvpair ...interface{}) {
	if h.opts.level.Enabled(LevelError) {
		Error(h).Print(kvpair...)
	}
}
