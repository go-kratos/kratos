package log

import (
	"context"
	"log"
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
	if err := c.logger.Log(level, kvs...); err != nil {
		return err
	}
	return nil
}

// MinLen At lest one kv pair
const MinLen = 2

// With logger fields.
func With(l Logger, kv ...interface{}) Logger {
	c, ok := l.(*logger)
	if !ok {
		return &logger{logger: l, prefix: kv, hasValuer: containsValuer(kv), ctx: context.Background()}
	}
	kvl := len(kv)
	// at lest one kv pair, or else return itself
	if kvl < MinLen {
		return l
	}
	// If len is an odd number, the last element is discard.
	if kvl&1 == 1 {
		kv = kv[:kvl-1]
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

// WithReplace : Add logger fields, and, if the key already exists, replace the older value with the new value.
func WithReplace(l Logger, kv ...interface{}) Logger {
	c, ok := l.(*logger)
	if !ok {
		return &logger{logger: l, prefix: kv, hasValuer: containsValuer(kv), ctx: context.Background()}
	}
	// at lest one kv pair, or else return itself
	if len(kv) < MinLen {
		return l
	}

	ca := len(c.prefix) + len(kv)
	kvs := make([]interface{}, 0, ca)
	filter := make(map[interface{}]bool, ca)

	// Add in reverse order, and the later will cover the earlier
	kvs = appendWithFilter(kv, kvs, filter)
	kvs = appendWithFilter(c.prefix, kvs, filter)
	// Reverse to asc order
	kvs = reverse(kvs)

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

func appendWithFilter(src []interface{}, to []interface{}, filter map[interface{}]bool) []interface{} {
	// If len is an odd number, the last element is discard.
	l := len(src)
	if l&1 == 1 {
		l--
	}
	for i := l - 1; i > 0; i = i - 2 {
		// exists key, skip kv
		if b, ok := filter[src[i-1]]; b && ok {
			i = i - 2
			continue
		}
		// append val
		to = append(to, src[i])
		// append key
		to = append(to, src[i-1])
		// save key to filter
		filter[src[i-1]] = true
	}
	return to
}

func reverse(slice []interface{}) []interface{} {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
