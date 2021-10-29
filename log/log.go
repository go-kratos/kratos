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
		l = lgr.WithContext(context.Background())
	}

	lgr := &logger{
		logs:      []Logger{l},
		prefix:    kv,
		hasValuer: containsValuer(kv),
		ctx:       context.Background(),
	}
	addSkipDepth(lgr, defaultDepth)
	return lgr
}

// WithContext returns a shallow copy of l with its context changed
// to ctx. The provided ctx must be non-nil.
func WithContext(ctx context.Context, l Logger) Logger {
	switch lgr := l.(type) {
	case *logger:
		if ctx == nil {
			return l
		}

		curDepth := getSkipDepth(lgr.ctx)
		return &logger{
			logs:      lgr.logs,
			prefix:    lgr.prefix,
			hasValuer: lgr.hasValuer,
			ctx:       setSkipDepth(ctx, curDepth),
		}
	case *Filter:
		l = lgr.WithContext(context.Background())
	}

	lgr := &logger{
		logs: []Logger{l},
		ctx:  ctx,
	}
	addSkipDepth(lgr, defaultDepth)
	return lgr
}

// MultiLogger wraps multi logger.
func MultiLogger(logs ...Logger) Logger {
	lgs := make([]Logger, 0, len(logs))
	for _, lgr := range logs {
		switch lg := lgr.(type) {
		case *logger:
			lgs = append(lgs, WithContext(context.Background(), lg))
		case *Filter:
			lgs = append(lgs, lg.WithContext(context.Background()))
		default:
			lgs = append(lgs, lg)
		}
	}

	mlg := WithContext(context.Background(), &logger{})
	addSkipDepth(mlg, 2)
	mLog := mlg.(*logger)
	mLog.logs = lgs
	addSkipDepth(mlg, 1)
	return mlg
}
