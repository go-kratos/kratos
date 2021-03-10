package log

import (
	"fmt"
)

// Helper is a logger helper.
type Helper struct {
	Logger
}

// NewHelper new a logger helper.
func NewHelper(name string, logger Logger) *Helper {
	return &Helper{
		Logger: With(logger, "module", name),
	}
}

// Debug logs a message at debug level.
func (h *Helper) Debug(a ...interface{}) {
	h.Print(LevelDebug, "message", fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func (h *Helper) Debugf(format string, a ...interface{}) {
	h.Print(LevelDebug, "message", fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func (h *Helper) Debugw(kvpair ...interface{}) {
	h.Print(LevelDebug, kvpair...)
}

// Info logs a message at info level.
func (h *Helper) Info(a ...interface{}) {
	h.Print(LevelInfo, "message", fmt.Sprint(a...))
}

// Infof logs a message at info level.
func (h *Helper) Infof(format string, a ...interface{}) {
	h.Print(LevelInfo, "message", fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func (h *Helper) Infow(kvpair ...interface{}) {
	h.Print(LevelInfo, kvpair...)
}

// Warn logs a message at warn level.
func (h *Helper) Warn(a ...interface{}) {
	h.Print(LevelWarn, "message", fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func (h *Helper) Warnf(format string, a ...interface{}) {
	h.Print(LevelWarn, "message", fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func (h *Helper) Warnw(kvpair ...interface{}) {
	h.Print(LevelWarn, kvpair...)
}

// Error logs a message at error level.
func (h *Helper) Error(a ...interface{}) {
	h.Print(LevelError, "message", fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func (h *Helper) Errorf(format string, a ...interface{}) {
	h.Print(LevelError, "message", fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func (h *Helper) Errorw(kvpair ...interface{}) {
	h.Print(LevelError, kvpair...)
}
