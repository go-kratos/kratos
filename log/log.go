package log

import (
	"context"
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

type logger struct {
	logs      []Logger
	prefix    []interface{}
	hasValuer bool
	ctx       context.Context
}

func (c *logger) Log(level Level, keyvals ...interface{}) error {
	kvs := make([]interface{}, 0, len(c.prefix)+len(keyvals))
	kvs = append(kvs, c.prefix...)
	if c.hasValuer {
		bindValues(c.ctx, kvs)
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
	if c, ok := l.(*logger); ok {
		kvs := make([]interface{}, 0, len(c.prefix)+len(kv))
		kvs = append(kvs, kv...)
		kvs = append(kvs, c.prefix...)
		return &logger{
			logs:      c.logs,
			prefix:    kvs,
			hasValuer: containsValuer(kvs),
			ctx:       c.ctx,
		}
	}
	return &logger{logs: []Logger{l}, prefix: kv, hasValuer: containsValuer(kv)}
}

// WithContext returns a shallow copy of l with its context changed
// to ctx. The provided ctx must be non-nil.
func WithContext(ctx context.Context, l Logger) Logger {
	if c, ok := l.(*logger); ok {
		return &logger{
			logs:      c.logs,
			prefix:    c.prefix,
			hasValuer: c.hasValuer,
			ctx:       ctx,
		}
	}
	return &logger{logs: []Logger{l}, ctx: ctx}
}

// MultiLogger wraps multi logger.
func MultiLogger(logs ...Logger) Logger {
	return &logger{logs: logs}
}
