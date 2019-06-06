package memcache

import (
	"context"

	"github.com/bilibili/kratos/pkg/container/pool"
	xtime "github.com/bilibili/kratos/pkg/time"
)

const (
	// Flag, 15(encoding) bit+ 17(compress) bit

	// FlagRAW default flag.
	FlagRAW = uint32(0)
	// FlagGOB gob encoding.
	FlagGOB = uint32(1) << 0
	// FlagJSON json encoding.
	FlagJSON = uint32(1) << 1
	// FlagProtobuf protobuf
	FlagProtobuf = uint32(1) << 2

	_flagEncoding = uint32(0xFFFF8000)

	// FlagGzip gzip compress.
	FlagGzip = uint32(1) << 15

	// left mv 31??? not work!!!
	flagLargeValue = uint32(1) << 30
)

// Item is an reply to be got or stored in a memcached server.
type Item struct {
	// Key is the Item's key (250 bytes maximum).
	Key string

	// Value is the Item's value.
	Value []byte

	// Object is the Item's object for use codec.
	Object interface{}

	// Flags are server-opaque flags whose semantics are entirely
	// up to the app.
	Flags uint32

	// Expiration is the cache expiration time, in seconds: either a relative
	// time from now (up to 1 month), or an absolute Unix epoch time.
	// Zero means the Item has no expiration time.
	Expiration int32

	// Compare and swap ID.
	cas uint64
}

// Conn represents a connection to a Memcache server.
// Command Reference: https://github.com/memcached/memcached/wiki/Commands
type Conn interface {
	// Close closes the connection.
	Close() error

	// Err returns a non-nil value if the connection is broken. The returned
	// value is either the first non-nil value returned from the underlying
	// network connection or a protocol parsing error. Applications should
	// close broken connections.
	Err() error

	// Add writes the given item, if no value already exists for its key.
	// ErrNotStored is returned if that condition is not met.
	Add(item *Item) error

	// Set writes the given item, unconditionally.
	Set(item *Item) error

	// Replace writes the given item, but only if the server *does* already
	// hold data for this key.
	Replace(item *Item) error

	// Get sends a command to the server for gets data.
	Get(key string) (*Item, error)

	// GetMulti is a batch version of Get. The returned map from keys to items
	// may have fewer elements than the input slice, due to memcache cache
	// misses. Each key must be at most 250 bytes in length.
	// If no error is returned, the returned map will also be non-nil.
	GetMulti(keys []string) (map[string]*Item, error)

	// Delete deletes the item with the provided key.
	// The error ErrNotFound is returned if the item didn't already exist in
	// the cache.
	Delete(key string) error

	// Increment atomically increments key by delta. The return value is the
	// new value after being incremented or an error. If the value didn't exist
	// in memcached the error is ErrNotFound. The value in memcached must be
	// an decimal number, or an error will be returned.
	// On 64-bit overflow, the new value wraps around.
	Increment(key string, delta uint64) (newValue uint64, err error)

	// Decrement atomically decrements key by delta. The return value is the
	// new value after being decremented or an error. If the value didn't exist
	// in memcached the error is ErrNotFound. The value in memcached must be
	// an decimal number, or an error will be returned. On underflow, the new
	// value is capped at zero and does not wrap around.
	Decrement(key string, delta uint64) (newValue uint64, err error)

	// CompareAndSwap writes the given item that was previously returned by
	// Get, if the value was neither modified or evicted between the Get and
	// the CompareAndSwap calls. The item's Key should not change between calls
	// but all other item fields may differ. ErrCASConflict is returned if the
	// value was modified in between the calls.
	// ErrNotStored is returned if the value was evicted in between the calls.
	CompareAndSwap(item *Item) error

	// Touch updates the expiry for the given key. The seconds parameter is
	// either a Unix timestamp or, if seconds is less than 1 month, the number
	// of seconds into the future at which time the item will expire.
	// ErrNotFound is returned if the key is not in the cache. The key must be
	// at most 250 bytes in length.
	Touch(key string, seconds int32) (err error)

	// Scan converts value read from the memcache into the following
	// common Go types and special types:
	//
	//    *string
	//    *[]byte
	//    *interface{}
	//
	Scan(item *Item, v interface{}) (err error)

	// Add writes the given item, if no value already exists for its key.
	// ErrNotStored is returned if that condition is not met.
	AddContext(ctx context.Context, item *Item) error

	// Set writes the given item, unconditionally.
	SetContext(ctx context.Context, item *Item) error

	// Replace writes the given item, but only if the server *does* already
	// hold data for this key.
	ReplaceContext(ctx context.Context, item *Item) error

	// Get sends a command to the server for gets data.
	GetContext(ctx context.Context, key string) (*Item, error)

	// GetMulti is a batch version of Get. The returned map from keys to items
	// may have fewer elements than the input slice, due to memcache cache
	// misses. Each key must be at most 250 bytes in length.
	// If no error is returned, the returned map will also be non-nil.
	GetMultiContext(ctx context.Context, keys []string) (map[string]*Item, error)

	// Delete deletes the item with the provided key.
	// The error ErrNotFound is returned if the item didn't already exist in
	// the cache.
	DeleteContext(ctx context.Context, key string) error

	// Increment atomically increments key by delta. The return value is the
	// new value after being incremented or an error. If the value didn't exist
	// in memcached the error is ErrNotFound. The value in memcached must be
	// an decimal number, or an error will be returned.
	// On 64-bit overflow, the new value wraps around.
	IncrementContext(ctx context.Context, key string, delta uint64) (newValue uint64, err error)

	// Decrement atomically decrements key by delta. The return value is the
	// new value after being decremented or an error. If the value didn't exist
	// in memcached the error is ErrNotFound. The value in memcached must be
	// an decimal number, or an error will be returned. On underflow, the new
	// value is capped at zero and does not wrap around.
	DecrementContext(ctx context.Context, key string, delta uint64) (newValue uint64, err error)

	// CompareAndSwap writes the given item that was previously returned by
	// Get, if the value was neither modified or evicted between the Get and
	// the CompareAndSwap calls. The item's Key should not change between calls
	// but all other item fields may differ. ErrCASConflict is returned if the
	// value was modified in between the calls.
	// ErrNotStored is returned if the value was evicted in between the calls.
	CompareAndSwapContext(ctx context.Context, item *Item) error

	// Touch updates the expiry for the given key. The seconds parameter is
	// either a Unix timestamp or, if seconds is less than 1 month, the number
	// of seconds into the future at which time the item will expire.
	// ErrNotFound is returned if the key is not in the cache. The key must be
	// at most 250 bytes in length.
	TouchContext(ctx context.Context, key string, seconds int32) (err error)
}

// Config memcache config.
type Config struct {
	*pool.Config

	Name         string // memcache name, for trace
	Proto        string
	Addr         string
	DialTimeout  xtime.Duration
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
}

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
func New(cfg *Config) *Memcache {
	return &Memcache{pool: NewPool(cfg)}
}

// Close close connection pool
func (mc *Memcache) Close() error {
	return mc.pool.Close()
}

// Conn direct get a connection
func (mc *Memcache) Conn(ctx context.Context) Conn {
	return mc.pool.Get(ctx)
}

// Set writes the given item, unconditionally.
func (mc *Memcache) Set(ctx context.Context, item *Item) (err error) {
	conn := mc.pool.Get(ctx)
	err = conn.SetContext(ctx, item)
	conn.Close()
	return
}

// Add writes the given item, if no value already exists for its key.
// ErrNotStored is returned if that condition is not met.
func (mc *Memcache) Add(ctx context.Context, item *Item) (err error) {
	conn := mc.pool.Get(ctx)
	err = conn.AddContext(ctx, item)
	conn.Close()
	return
}

// Replace writes the given item, but only if the server *does* already hold data for this key.
func (mc *Memcache) Replace(ctx context.Context, item *Item) (err error) {
	conn := mc.pool.Get(ctx)
	err = conn.ReplaceContext(ctx, item)
	conn.Close()
	return
}

// CompareAndSwap writes the given item that was previously returned by Get
func (mc *Memcache) CompareAndSwap(ctx context.Context, item *Item) (err error) {
	conn := mc.pool.Get(ctx)
	err = conn.CompareAndSwapContext(ctx, item)
	conn.Close()
	return
}

// Get sends a command to the server for gets data.
func (mc *Memcache) Get(ctx context.Context, key string) *Reply {
	conn := mc.pool.Get(ctx)
	item, err := conn.GetContext(ctx, key)
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
func (mc *Memcache) GetMulti(ctx context.Context, keys []string) (*Replies, error) {
	conn := mc.pool.Get(ctx)
	items, err := conn.GetMultiContext(ctx, keys)
	rs := &Replies{err: err, items: items, conn: conn, usedItems: make(map[string]struct{}, len(keys))}
	if (err != nil) || (len(items) == 0) {
		rs.Close()
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
		rs.Close()
		return ErrNotFound
	}
	rs.usedItems[key] = struct{}{}
	err = rs.conn.Scan(item, v)
	if (err != nil) || (len(rs.items) == len(rs.usedItems)) {
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
func (mc *Memcache) Touch(ctx context.Context, key string, timeout int32) (err error) {
	conn := mc.pool.Get(ctx)
	err = conn.TouchContext(ctx, key, timeout)
	conn.Close()
	return
}

// Delete deletes the item with the provided key.
func (mc *Memcache) Delete(ctx context.Context, key string) (err error) {
	conn := mc.pool.Get(ctx)
	err = conn.DeleteContext(ctx, key)
	conn.Close()
	return
}

// Increment atomically increments key by delta.
func (mc *Memcache) Increment(ctx context.Context, key string, delta uint64) (newValue uint64, err error) {
	conn := mc.pool.Get(ctx)
	newValue, err = conn.IncrementContext(ctx, key, delta)
	conn.Close()
	return
}

// Decrement atomically decrements key by delta.
func (mc *Memcache) Decrement(ctx context.Context, key string, delta uint64) (newValue uint64, err error) {
	conn := mc.pool.Get(ctx)
	newValue, err = conn.DecrementContext(ctx, key, delta)
	conn.Close()
	return
}
