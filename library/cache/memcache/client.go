package memcache

import (
	"context"
)

// Memcache memcache client
type Memcache struct {
	pool *Pool
}

// Reply is the result of Get
type Reply struct {
	err    error
	item   *Item
	conn   Conn
	closed bool
}

// Replies is the result of GetMulti
type Replies struct {
	err       error
	items     map[string]*Item
	usedItems map[string]struct{}
	conn      Conn
	closed    bool
}

// New get a memcache client
func New(c *Config) *Memcache {
	return &Memcache{pool: NewPool(c)}
}

// Close close connection pool
func (mc *Memcache) Close() error {
	return mc.pool.Close()
}

// Conn direct get a connection
func (mc *Memcache) Conn(c context.Context) Conn {
	return mc.pool.Get(c)
}

// Set writes the given item, unconditionally.
func (mc *Memcache) Set(c context.Context, item *Item) (err error) {
	conn := mc.pool.Get(c)
	err = conn.Set(item)
	conn.Close()
	return
}

// Add writes the given item, if no value already exists for its key.
// ErrNotStored is returned if that condition is not met.
func (mc *Memcache) Add(c context.Context, item *Item) (err error) {
	conn := mc.pool.Get(c)
	err = conn.Add(item)
	conn.Close()
	return
}

// Replace writes the given item, but only if the server *does* already hold data for this key.
func (mc *Memcache) Replace(c context.Context, item *Item) (err error) {
	conn := mc.pool.Get(c)
	err = conn.Replace(item)
	conn.Close()
	return
}

// CompareAndSwap writes the given item that was previously returned by Get
func (mc *Memcache) CompareAndSwap(c context.Context, item *Item) (err error) {
	conn := mc.pool.Get(c)
	err = conn.CompareAndSwap(item)
	conn.Close()
	return
}

// Get sends a command to the server for gets data.
func (mc *Memcache) Get(c context.Context, key string) *Reply {
	conn := mc.pool.Get(c)
	item, err := conn.Get(key)
	if err != nil {
		conn.Close()
	}
	return &Reply{err: err, item: item, conn: conn}
}

// Item get raw Item
func (r *Reply) Item() *Item {
	return r.item
}

// Scan converts value, read from the memcache
func (r *Reply) Scan(v interface{}) (err error) {
	if r.err != nil {
		return r.err
	}
	err = r.conn.Scan(r.item, v)
	if !r.closed {
		r.conn.Close()
		r.closed = true
	}
	return
}

// GetMulti is a batch version of Get
func (mc *Memcache) GetMulti(c context.Context, keys []string) (*Replies, error) {
	conn := mc.pool.Get(c)
	items, err := conn.GetMulti(keys)
	rs := &Replies{err: err, items: items, conn: conn, usedItems: make(map[string]struct{}, len(keys))}
	if err != nil {
		conn.Close()
		rs.closed = true
	}
	return rs, err
}

// Close close rows.
func (rs *Replies) Close() (err error) {
	if !rs.closed {
		err = rs.conn.Close()
		rs.closed = true
	}
	return
}

// Item get Item from rows
func (rs *Replies) Item(key string) *Item {
	return rs.items[key]
}

// Scan converts value, read from key in rows
func (rs *Replies) Scan(key string, v interface{}) (err error) {
	if rs.err != nil {
		return rs.err
	}
	item, ok := rs.items[key]
	if !ok {
		return ErrNotFound
	}
	rs.usedItems[key] = struct{}{}
	err = rs.conn.Scan(item, v)
	shouldClose := len(rs.items) == len(rs.usedItems)
	if shouldClose {
		rs.Close()
	}
	return
}

// Keys keys of result
func (rs *Replies) Keys() (keys []string) {
	keys = make([]string, 0, len(rs.items))
	for key := range rs.items {
		keys = append(keys, key)
	}
	return
}

// Touch updates the expiry for the given key.
func (mc *Memcache) Touch(c context.Context, key string, timeout int32) (err error) {
	conn := mc.pool.Get(c)
	err = conn.Touch(key, timeout)
	conn.Close()
	return
}

// Delete deletes the item with the provided key.
func (mc *Memcache) Delete(c context.Context, key string) (err error) {
	conn := mc.pool.Get(c)
	err = conn.Delete(key)
	conn.Close()
	return
}

// Increment atomically increments key by delta.
func (mc *Memcache) Increment(c context.Context, key string, delta uint64) (newValue uint64, err error) {
	conn := mc.pool.Get(c)
	newValue, err = conn.Increment(key, delta)
	conn.Close()
	return
}

// Decrement atomically decrements key by delta.
func (mc *Memcache) Decrement(c context.Context, key string, delta uint64) (newValue uint64, err error) {
	conn := mc.pool.Get(c)
	newValue, err = conn.Decrement(key, delta)
	conn.Close()
	return
}
