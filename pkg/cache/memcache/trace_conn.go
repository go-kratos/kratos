package memcache

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/net/trace"
)

const (
	_slowLogDuration = time.Millisecond * 250
)

func newTraceConn(conn Conn, address string) Conn {
	tags := []trace.Tag{
		trace.String(trace.TagSpanKind, "client"),
		trace.String(trace.TagComponent, "cache/memcache"),
		trace.String(trace.TagPeerService, "memcache"),
		trace.String(trace.TagPeerAddress, address),
	}
	return &traceConn{Conn: conn, tags: tags}
}

type traceConn struct {
	Conn
	tags []trace.Tag
}

func (t *traceConn) setTrace(ctx context.Context, action, statement string) func(error) error {
	now := time.Now()
	parent, ok := trace.FromContext(ctx)
	if !ok {
		return func(err error) error { return err }
	}
	span := parent.Fork("", "Memcache:"+action)
	span.SetTag(t.tags...)
	span.SetTag(trace.String(trace.TagDBStatement, action+" "+statement))
	return func(err error) error {
		span.Finish(&err)
		t := time.Since(now)
		if t > _slowLogDuration {
			log.Warn("memcache slow log action: %s key: %s time: %v", action, statement, t)
		}
		return err
	}
}

func (t *traceConn) AddContext(ctx context.Context, item *Item) error {
	finishFn := t.setTrace(ctx, "Add", item.Key)
	return finishFn(t.Conn.Add(item))
}

func (t *traceConn) SetContext(ctx context.Context, item *Item) error {
	finishFn := t.setTrace(ctx, "Set", item.Key)
	return finishFn(t.Conn.Set(item))
}

func (t *traceConn) ReplaceContext(ctx context.Context, item *Item) error {
	finishFn := t.setTrace(ctx, "Replace", item.Key)
	return finishFn(t.Conn.Replace(item))
}

func (t *traceConn) GetContext(ctx context.Context, key string) (*Item, error) {
	finishFn := t.setTrace(ctx, "Get", key)
	item, err := t.Conn.Get(key)
	return item, finishFn(err)
}

func (t *traceConn) GetMultiContext(ctx context.Context, keys []string) (map[string]*Item, error) {
	finishFn := t.setTrace(ctx, "GetMulti", strings.Join(keys, " "))
	items, err := t.Conn.GetMulti(keys)
	return items, finishFn(err)
}

func (t *traceConn) DeleteContext(ctx context.Context, key string) error {
	finishFn := t.setTrace(ctx, "Delete", key)
	return finishFn(t.Conn.Delete(key))
}

func (t *traceConn) IncrementContext(ctx context.Context, key string, delta uint64) (newValue uint64, err error) {
	finishFn := t.setTrace(ctx, "Increment", key+" "+strconv.FormatUint(delta, 10))
	newValue, err = t.Conn.Increment(key, delta)
	return newValue, finishFn(err)
}

func (t *traceConn) DecrementContext(ctx context.Context, key string, delta uint64) (newValue uint64, err error) {
	finishFn := t.setTrace(ctx, "Decrement", key+" "+strconv.FormatUint(delta, 10))
	newValue, err = t.Conn.Decrement(key, delta)
	return newValue, finishFn(err)
}

func (t *traceConn) CompareAndSwapContext(ctx context.Context, item *Item) error {
	finishFn := t.setTrace(ctx, "CompareAndSwap", item.Key)
	return finishFn(t.Conn.CompareAndSwap(item))
}

func (t *traceConn) TouchContext(ctx context.Context, key string, seconds int32) (err error) {
	finishFn := t.setTrace(ctx, "Touch", key+" "+strconv.Itoa(int(seconds)))
	return finishFn(t.Conn.Touch(key, seconds))
}
