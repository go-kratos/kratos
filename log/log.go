package log

import (
	"context"
	"log"
)

// DefaultLogger is default logger.
var DefaultLogger Logger = NewStdLogger(log.Writer())

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
	switch lgr := l.(type) {
	case *logger:
		kvs := make([]interface{}, 0, len(lgr.prefix)+len(kv))
		kvs = append(kvs, kv...)
		kvs = append(kvs, lgr.prefix...)
		return &logger{
			logs:      lgr.logs,
			prefix:    kvs,
			hasValuer: containsValuer(kvs),
			ctx:       lgr.ctx,
		}
	case *Filter:
		return &Filter{
			logger: With(lgr.logger, kv...),
			level:  lgr.level,
			key:    lgr.key,
			filter: lgr.filter,
		}
	case *Helper:
		return &Helper{
			logger: With(lgr.logger, kv...),
			msgKey: lgr.msgKey,
		}
	default:
		lgr = &logger{
			logs:      []Logger{l},
			prefix:    kv,
			hasValuer: containsValuer(kv),
		}
		return WithContext(context.Background(), lgr)
	}
}

// WithContext returns a shallow copy of l with its context changed
// to ctx. The provided ctx must be non-nil.
func WithContext(ctx context.Context, l Logger) Logger {
	if ctx == nil {
		return l
	}
	return withContext(ctx, l, 1)
}

func withContext(ctx context.Context, l Logger, depth int) Logger {
	switch lgr := l.(type) {
	case *logger:
		lgs := make([]Logger, 0, len(lgr.logs))
		for _, subLog := range lgr.logs {
			lgs = append(lgs, withContext(ctx, subLog, depth+1))
		}
		return &logger{
			logs:      lgs,
			prefix:    lgr.prefix,
			hasValuer: lgr.hasValuer,
			ctx:       setSkipDepth(ctx, depth),
		}
	case *Filter:
		return &Filter{
			logger: withContext(ctx, lgr.logger, depth+1),
			level:  lgr.level,
			key:    lgr.key,
			filter: lgr.filter,
		}
	case *Helper:
		return &Helper{
			logger: withContext(ctx, lgr.logger, depth+1),
			msgKey: lgr.msgKey,
		}
	default:
		return l // Other log struct cannot be bound to a context
	}
}

// MultiLogger wraps multi logger.
func MultiLogger(logs ...Logger) Logger {
	mlg := &logger{
		logs: logs,
	}
	return WithContext(context.Background(), mlg)
}
