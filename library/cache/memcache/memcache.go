package memcache

import (
	"context"
)

// Error represents an error returned in a command reply.
type Error string

func (err Error) Error() string { return string(err) }

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
	// The error ErrCacheMiss is returned if the item didn't already exist in
	// the cache.
	Delete(key string) error

	// Increment atomically increments key by delta. The return value is the
	// new value after being incremented or an error. If the value didn't exist
	// in memcached the error is ErrCacheMiss. The value in memcached must be
	// an decimal number, or an error will be returned.
	// On 64-bit overflow, the new value wraps around.
	Increment(key string, delta uint64) (newValue uint64, err error)

	// Decrement atomically decrements key by delta. The return value is the
	// new value after being decremented or an error. If the value didn't exist
	// in memcached the error is ErrCacheMiss. The value in memcached must be
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
	//ErrCacheMiss is returned if the key is not in the cache. The key must be
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

	// WithContext return a Conn with its context changed to ctx
	// the context controls the entire lifetime of Conn before you change it
	// NOTE: this method is not thread-safe
	WithContext(ctx context.Context) Conn
}
