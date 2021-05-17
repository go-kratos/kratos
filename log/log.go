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
	Log(level Level, keyvals ...interface{}) error
}

type context struct {
	logs      []Logger
	prefix    []interface{}
	hasValuer bool
}

func (c *context) Log(level Level, keyvals ...interface{}) error {
	kvs := make([]interface{}, 0, len(c.prefix)+len(keyvals))
	kvs = append(kvs, c.prefix...)
	if c.hasValuer {
		bindValues(kvs)
	}
	kvs = append(kvs, keyvals...)
	for _, l := range c.logs {
		if err := l.Log(level, kvs...); err != nil {
			return err
		}
	}
	return nil
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

// MultiLogger wraps multi logger.
func MultiLogger(logs ...Logger) Logger {
	return &context{logs: logs}
}
