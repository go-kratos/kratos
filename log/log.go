package log

import (
	"log"
)

var (
	// DefaultLogger is default logger.
	DefaultLogger Logger = NewStdLogger(log.Writer())
)

// Logger is a logger interface.
type Logger interface {
	Print(pairs ...interface{})
}

type context struct {
	logs   []Logger
	prefix []interface{}
}

func (c *context) Print(a ...interface{}) {
	kvs := make([]interface{}, 0, len(c.prefix)+len(a))
	kvs = append(kvs, c.prefix...)
	kvs = append(kvs, a...)
	for _, log := range c.logs {
		log.Print(kvs...)
	}
}

// With with logger fields.
func With(l Logger, a ...interface{}) Logger {
	if c, ok := l.(*context); ok {
		kvs := make([]interface{}, 0, len(c.prefix)+len(a))
		kvs = append(kvs, a...)
		kvs = append(kvs, c.prefix...)
		return &context{
			logs:   c.logs,
			prefix: kvs,
		}
	}
	return &context{logs: []Logger{l}, prefix: a}
}

// Wrap wraps multi logger.
func Wrap(logs ...Logger) Logger {
	return &context{logs: logs}
}

// Debug returns a debug logger.
func Debug(log Logger) Logger {
	return With(log, LevelKey, LevelDebug)
}

// Info returns a info logger.
func Info(log Logger) Logger {
	return With(log, LevelKey, LevelInfo)
}

// Warn return a warn logger.
func Warn(log Logger) Logger {
	return With(log, LevelKey, LevelWarn)
}

// Error returns a error logger.
func Error(log Logger) Logger {
	return With(log, LevelKey, LevelError)
}
