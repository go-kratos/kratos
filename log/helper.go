package log

import (
	"fmt"
)

// Helper is a logger helper.
type Helper struct {
	debug Logger
	info  Logger
	warn  Logger
	err   Logger
}

// NewHelper new a logger helper.
func NewHelper(name string, logger Logger) *Helper {
	log := With(logger, "module", name)
	return &Helper{
		debug: Debug(log),
		info:  Info(log),
		warn:  Warn(log),
		err:   Error(log),
	}
}

// Debug logs a message at debug level.
func (h *Helper) Debug(a ...interface{}) {
	h.debug.Print("message", fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func (h *Helper) Debugf(format string, a ...interface{}) {
	h.debug.Print("message", fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func (h *Helper) Debugw(pairs ...interface{}) {
	h.debug.Print(pairs...)
}

// Info logs a message at info level.
func (h *Helper) Info(a ...interface{}) {
	h.info.Print("message", fmt.Sprint(a...))
}

// Infof logs a message at info level.
func (h *Helper) Infof(format string, a ...interface{}) {
	h.info.Print("message", fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func (h *Helper) Infow(pairs ...interface{}) {
	h.info.Print(pairs...)
}

// Warn logs a message at warn level.
func (h *Helper) Warn(a ...interface{}) {
	h.warn.Print("message", fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func (h *Helper) Warnf(format string, a ...interface{}) {
	h.warn.Print("message", fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func (h *Helper) Warnw(pairs ...interface{}) {
	h.warn.Print(pairs...)
}

// Error logs a message at error level.
func (h *Helper) Error(a ...interface{}) {
	h.err.Print("message", fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func (h *Helper) Errorf(format string, a ...interface{}) {
	h.err.Print("message", fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func (h *Helper) Errorw(pairs ...interface{}) {
	h.err.Print(pairs...)
}
