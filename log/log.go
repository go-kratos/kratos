package log

import (
	"context"
	"log"
	"reflect"
)

// DefaultLogger is default logger.
var DefaultLogger = NewStdLogger(log.Writer())

// Logger is a logger interface.
type Logger interface {
	Log(level Level, keyvals ...interface{}) error
}

type logger struct {
	logger    Logger
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
	if err := c.logger.Log(level, uniqKeys(kvs)...); err != nil {
		return err
	}
	return nil
}

// uniqKeys Under the same key, the Val corresponding to the last key shall prevail
func uniqKeys(kvs []interface{}) []interface{} {
	if len(kvs) == 0 {
		return kvs
	}
	if len(kvs)&1 == 1 {
		kvs = append(kvs, "KEYVALS UNPAIRED") // key missing val, add a virtual value
	}
	uniqMap := make(map[reflect.Value]interface{})
	keys := make([]interface{}, 0, len(kvs))
	for i := 0; i < len(kvs); i = i + 2 {
		key, val := kvs[i], kvs[i+1]
		rvKey := reflect.ValueOf(key) // key missing val, map cannot be stored log.Valuer, so reflect.Value is used as the key here
		if _, ok := uniqMap[rvKey]; !ok {
			keys = append(keys, key) // ensure the data sequence is consistent
		}
		uniqMap[rvKey] = val // ensure that the val corresponding to the key is the latest
	}
	uniq := make([]interface{}, 0, 2*len(keys))
	for _, key := range keys {
		uniq = append(uniq, key, uniqMap[reflect.ValueOf(key)])
	}
	return uniq
}

// With with logger fields.
func With(l Logger, kv ...interface{}) Logger {
	c, ok := l.(*logger)
	if !ok {
		return &logger{logger: l, prefix: kv, hasValuer: containsValuer(kv), ctx: context.Background()}
	}
	kvs := make([]interface{}, 0, len(c.prefix)+len(kv))
	kvs = append(kvs, c.prefix...)
	kvs = append(kvs, kv...)
	return &logger{
		logger:    c.logger,
		prefix:    kvs,
		hasValuer: containsValuer(kvs),
		ctx:       c.ctx,
	}
}

// WithContext returns a shallow copy of l with its context changed
// to ctx. The provided ctx must be non-nil.
func WithContext(ctx context.Context, l Logger) Logger {
	c, ok := l.(*logger)
	if !ok {
		return &logger{logger: l, ctx: ctx}
	}
	return &logger{
		logger:    c.logger,
		prefix:    c.prefix,
		hasValuer: c.hasValuer,
		ctx:       ctx,
	}
}
