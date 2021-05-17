package log

import (
	"fmt"
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

// Log .
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
