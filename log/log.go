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
	Print(kv ...interface{})
}

type context struct {
	logs      []Logger
	prefix    []interface{}
	hasValuer bool
}

func (c *context) Print(kv ...interface{}) {
	kvs := make([]interface{}, 0, len(c.prefix)+len(kv))
	kvs = append(kvs, c.prefix...)
	if c.hasValuer {
		bindValues(kvs)
	}
	kvs = append(kvs, kv...)
	for _, l := range c.logs {
		l.Print(kvs...)
	}
}

// With with logger fields.
func With(l Logger, kv ...interface{}) Logger {
	if c, ok := l.(*context); ok {
		kvs := make([]interface{}, 0, len(c.prefix)+len(kv))
		kvs = append(kvs, kv...)
		kvs = append(kvs, c.prefix...)
		return &context{
			logs:      c.logs,
			prefix:    kvs,
			hasValuer: containsValuer(kvs),
		}
	}
	return &context{logs: []Logger{l}, prefix: kv, hasValuer: containsValuer(kv)}
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
