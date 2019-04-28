package metadata

import (
	"context"
	"fmt"
	"strconv"
)

// MD is a mapping from metadata keys to values.
type MD map[string]interface{}

type mdKey struct{}

// Len returns the number of items in md.
func (md MD) Len() int {
	return len(md)
}

// Copy returns a copy of md.
func (md MD) Copy() MD {
	return Join(md)
}

// New creates an MD from a given key-value map.
func New(m map[string]interface{}) MD {
	md := MD{}
	for k, val := range m {
		md[k] = val
	}
	return md
}

// Join joins any number of mds into a single MD.
// The order of values for each key is determined by the order in which
// the mds containing those values are presented to Join.
func Join(mds ...MD) MD {
	out := MD{}
	for _, md := range mds {
		for k, v := range md {
			out[k] = v
		}
	}
	return out
}

// Pairs returns an MD formed by the mapping of key, value ...
// Pairs panics if len(kv) is odd.
func Pairs(kv ...interface{}) MD {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: Pairs got the odd number of input pairs for metadata: %d", len(kv)))
	}
	md := MD{}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = s.(string)
			continue
		}
		md[key] = s
	}
	return md
}

// NewContext creates a new context with md attached.
func NewContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, mdKey{}, md)
}

// FromContext returns the incoming metadata in ctx if it exists.  The
// returned MD should not be modified. Writing to it may cause races.
// Modification should be made to copies of the returned MD.
func FromContext(ctx context.Context) (md MD, ok bool) {
	md, ok = ctx.Value(mdKey{}).(MD)
	return
}

// String get string value from metadata in context
func String(ctx context.Context, key string) string {
	md, ok := ctx.Value(mdKey{}).(MD)
	if !ok {
		return ""
	}
	str, _ := md[key].(string)
	return str
}

// Int64 get int64 value from metadata in context
func Int64(ctx context.Context, key string) int64 {
	md, ok := ctx.Value(mdKey{}).(MD)
	if !ok {
		return 0
	}
	i64, _ := md[key].(int64)
	return i64
}

// Value get value from metadata in context return nil if not found
func Value(ctx context.Context, key string) interface{} {
	md, ok := ctx.Value(mdKey{}).(MD)
	if !ok {
		return nil
	}
	return md[key]
}

// WithContext return no deadline context and retain metadata.
func WithContext(c context.Context) context.Context {
	md, ok := FromContext(c)
	if ok {
		nmd := md.Copy()
		// NOTE: temporary delete prevent asynchronous task reuse finished task
		delete(nmd, Trace)
		return NewContext(context.Background(), nmd)
	}
	return context.Background()
}

// Bool get boolean from metadata in context use strconv.Parse.
func Bool(ctx context.Context, key string) bool {
	md, ok := ctx.Value(mdKey{}).(MD)
	if !ok {
		return false
	}

	switch md[key].(type) {
	case bool:
		return md[key].(bool)
	case string:
		ok, _ = strconv.ParseBool(md[key].(string))
		return ok
	default:
		return false
	}
}
