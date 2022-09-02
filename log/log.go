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

// With logger fields.
func With(l Logger, kv ...interface{}) Logger {
	c, ok := l.(*logger)
	if !ok {
		return &logger{logger: l, prefix: kv, hasValuer: containsValuer(kv), ctx: context.Background()}
	}
	ca := len(c.prefix) + len(kv)
	kvs := make([]interface{}, 0, ca)
	m := make(map[interface{}]bool, ca)

	// Add in reverse order, and the later will cover the earlier
	kvs = appendWithFilter(kv, kvs, m)
	kvs = appendWithFilter(c.prefix, kvs, m)
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
	vl := len(src)
	if vl&1 == 1 {
		vl--
	}
	for i := vl - 1; i > 0; i = i - 2 {
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
