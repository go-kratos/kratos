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
	log    Logger
	prefix []interface{}
	suffix []interface{}
}

func (c *context) Print(a ...interface{}) {
	kvs := make([]interface{}, 0, len(c.prefix)+len(c.suffix)+len(a))
	kvs = append(kvs, c.prefix...)
	kvs = append(kvs, a...)
	kvs = append(kvs, c.suffix...)
	c.log.Print(kvs...)
}

// With with logger suffix.
func With(l Logger, a ...interface{}) Logger {
	if c, ok := l.(*context); ok {
		kvs := make([]interface{}, 0, len(c.suffix)+len(a))
		kvs = append(kvs, c.suffix...)
		kvs = append(kvs, a...)
		return &context{
			log:    c.log,
			prefix: c.prefix,
			suffix: kvs,
		}
	}
	return &context{log: l, suffix: a}
}

// Prefix with logger prefix.
func Prefix(l Logger, a ...interface{}) Logger {
	if c, ok := l.(*context); ok {
		kvs := make([]interface{}, 0, len(c.prefix)+len(a))
		kvs = append(kvs, a...)
		kvs = append(kvs, c.prefix...)
		return &context{
			log:    c.log,
			prefix: kvs,
			suffix: c.suffix,
		}
	}
	return &context{log: l, prefix: a}
}

// Debug returns a debug logger.
func Debug(log Logger) Logger {
	return Prefix(log, LevelKey, LevelDebug)
}

// Info returns a info logger.
func Info(log Logger) Logger {
	return Prefix(log, LevelKey, LevelInfo)
}

// Warn return a warn logger.
func Warn(log Logger) Logger {
	return Prefix(log, LevelKey, LevelWarn)
}

// Error returns a error logger.
func Error(log Logger) Logger {
	return Prefix(log, LevelKey, LevelError)
}
