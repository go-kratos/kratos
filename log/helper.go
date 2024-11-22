package log

import (
	"context"
	"fmt"
	"os"
)

// DefaultMessageKey default message key.
var DefaultMessageKey = "msg"

// Option is Helper option.
type Option func(*Helper)

// Helper is a logger helper.
type Helper struct {
	logger  Logger
	msgKey  string
	sprint  func(...interface{}) string
	sprintf func(format string, a ...interface{}) string
}

// WithMessageKey with message key.
func WithMessageKey(k string) Option {
	return func(opts *Helper) {
		opts.msgKey = k
	}
}

// WithSprint with sprint
func WithSprint(sprint func(...interface{}) string) Option {
	return func(opts *Helper) {
		opts.sprint = sprint
	}
}

// WithSprintf with sprintf
func WithSprintf(sprintf func(format string, a ...interface{}) string) Option {
	return func(opts *Helper) {
		opts.sprintf = sprintf
	}
}

// NewHelper new a logger helper.
func NewHelper(logger Logger, opts ...Option) *Helper {
	options := &Helper{
		msgKey:  DefaultMessageKey, // default message key
		logger:  logger,
		sprint:  fmt.Sprint,
		sprintf: fmt.Sprintf,
	}
	for _, o := range opts {
		o(options)
	}
	return options
}

// WithContext returns a shallow copy of h with its context changed
// to ctx. The provided ctx must be non-nil.
func (h *Helper) WithContext(ctx context.Context) *Helper {
	return &Helper{
		msgKey:  h.msgKey,
		logger:  WithContext(ctx, h.logger),
		sprint:  h.sprint,
		sprintf: h.sprintf,
	}
}

// Enabled returns true if the given level above this level.
// It delegates to the underlying *Filter.
func (h *Helper) Enabled(level Level) bool {
	if l, ok := h.logger.(*Filter); ok {
		return level >= l.level
	}
	return true
}

// Logger returns logger in the helper.
func (h *Helper) Logger() Logger {
	return h.logger
}

// Log Print log by level and keyvals.
func (h *Helper) Log(level Level, keyvals ...interface{}) {
	_ = h.logger.Log(level, keyvals...)
}

// Debug logs a message at debug level.
func (h *Helper) Debug(a ...interface{}) {
	if !h.Enabled(LevelDebug) {
		return
	}
	_ = h.logger.Log(LevelDebug, h.msgKey, h.sprint(a...))
}

// Debugf logs a message at debug level.
func (h *Helper) Debugf(format string, a ...interface{}) {
	if !h.Enabled(LevelDebug) {
		return
	}
	_ = h.logger.Log(LevelDebug, h.msgKey, h.sprintf(format, a...))
}

// Debugw logs a message at debug level.
func (h *Helper) Debugw(keyvals ...interface{}) {
	_ = h.logger.Log(LevelDebug, keyvals...)
}

// Info logs a message at info level.
func (h *Helper) Info(a ...interface{}) {
	if !h.Enabled(LevelInfo) {
		return
	}
	_ = h.logger.Log(LevelInfo, h.msgKey, h.sprint(a...))
}

// Infof logs a message at info level.
func (h *Helper) Infof(format string, a ...interface{}) {
	if !h.Enabled(LevelInfo) {
		return
	}
	_ = h.logger.Log(LevelInfo, h.msgKey, h.sprintf(format, a...))
}

// Infow logs a message at info level.
func (h *Helper) Infow(keyvals ...interface{}) {
	_ = h.logger.Log(LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func (h *Helper) Warn(a ...interface{}) {
	if !h.Enabled(LevelWarn) {
		return
	}
	_ = h.logger.Log(LevelWarn, h.msgKey, h.sprint(a...))
}

// Warnf logs a message at warnf level.
func (h *Helper) Warnf(format string, a ...interface{}) {
	if !h.Enabled(LevelWarn) {
		return
	}
	_ = h.logger.Log(LevelWarn, h.msgKey, h.sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func (h *Helper) Warnw(keyvals ...interface{}) {
	_ = h.logger.Log(LevelWarn, keyvals...)
}

// Error logs a message at error level.
func (h *Helper) Error(a ...interface{}) {
	if !h.Enabled(LevelError) {
		return
	}
	_ = h.logger.Log(LevelError, h.msgKey, h.sprint(a...))
}

// Errorf logs a message at error level.
func (h *Helper) Errorf(format string, a ...interface{}) {
	if !h.Enabled(LevelError) {
		return
	}
	_ = h.logger.Log(LevelError, h.msgKey, h.sprintf(format, a...))
}

// Errorw logs a message at error level.
func (h *Helper) Errorw(keyvals ...interface{}) {
	_ = h.logger.Log(LevelError, keyvals...)
}

// Fatal logs a message at fatal level.
func (h *Helper) Fatal(a ...interface{}) {
	_ = h.logger.Log(LevelFatal, h.msgKey, h.sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func (h *Helper) Fatalf(format string, a ...interface{}) {
	_ = h.logger.Log(LevelFatal, h.msgKey, h.sprintf(format, a...))
	os.Exit(1)
}

// Fatalw logs a message at fatal level.
func (h *Helper) Fatalw(keyvals ...interface{}) {
	_ = h.logger.Log(LevelFatal, keyvals...)
	os.Exit(1)
}
