package log

import (
	"context"
	"fmt"
	"os"
)

// Helper is a logger helper.
type Helper struct {
	logger Logger
}

// NewHelper new a logger helper.
func NewHelper(logger Logger) *Helper {
	return &Helper{
		logger: logger,
	}
}

// WithContext returns a shallow copy of h with its context changed
// to ctx. The provided ctx must be non-nil.
func (h *Helper) WithContext(ctx context.Context) *Helper {
	return &Helper{
		logger: WithContext(ctx, h.logger),
	}
}

// Log Print log by level and keyvals.
func (h *Helper) Log(level Level, keyvals ...interface{}) {
	h.logger.Log(level, keyvals...)
}

// Debug logs a message at debug level.
func (h *Helper) Debug(a ...interface{}) {
	h.logger.Log(LevelDebug, "msg", fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func (h *Helper) Debugf(format string, a ...interface{}) {
	h.logger.Log(LevelDebug, "msg", fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func (h *Helper) Debugw(keyvals ...interface{}) {
	h.logger.Log(LevelDebug, keyvals...)
}

// Info logs a message at info level.
func (h *Helper) Info(a ...interface{}) {
	h.logger.Log(LevelInfo, "msg", fmt.Sprint(a...))
}

// Infof logs a message at info level.
func (h *Helper) Infof(format string, a ...interface{}) {
	h.logger.Log(LevelInfo, "msg", fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func (h *Helper) Infow(keyvals ...interface{}) {
	h.logger.Log(LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func (h *Helper) Warn(a ...interface{}) {
	h.logger.Log(LevelWarn, "msg", fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func (h *Helper) Warnf(format string, a ...interface{}) {
	h.logger.Log(LevelWarn, "msg", fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func (h *Helper) Warnw(keyvals ...interface{}) {
	h.logger.Log(LevelWarn, keyvals...)
}

// Error logs a message at error level.
func (h *Helper) Error(a ...interface{}) {
	h.logger.Log(LevelError, "msg", fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func (h *Helper) Errorf(format string, a ...interface{}) {
	h.logger.Log(LevelError, "msg", fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func (h *Helper) Errorw(keyvals ...interface{}) {
	h.logger.Log(LevelError, keyvals...)
}

// Fatal logs a message at fatal level.
func (h *Helper) Fatal(a ...interface{}) {
	h.logger.Log(LevelFatal, "msg", fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func (h *Helper) Fatalf(format string, a ...interface{}) {
	h.logger.Log(LevelFatal, "msg", fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Fatalw logs a message at fatal level.
func (h *Helper) Fatalw(keyvals ...interface{}) {
	h.logger.Log(LevelFatal, keyvals...)
	os.Exit(1)
}
