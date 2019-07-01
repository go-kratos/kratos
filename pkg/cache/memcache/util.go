package memcache

import (
	"context"
	"time"

	"github.com/gogo/protobuf/proto"
)

func legalKey(key string) bool {
	if len(key) > 250 || len(key) == 0 {
		return false
	}
	for i := 0; i < len(key); i++ {
		if key[i] <= ' ' || key[i] == 0x7f {
			return false
		}
	}
	return true
}

// MockWith error
func MockWith(err error) Conn {
	return errConn{err}
}

type errConn struct{ err error }

func (c errConn) Err() error                                                       { return c.err }
func (c errConn) Close() error                                                     { return c.err }
func (c errConn) Add(*Item) error                                                  { return c.err }
func (c errConn) Set(*Item) error                                                  { return c.err }
func (c errConn) Replace(*Item) error                                              { return c.err }
func (c errConn) CompareAndSwap(*Item) error                                       { return c.err }
func (c errConn) Get(string) (*Item, error)                                        { return nil, c.err }
func (c errConn) GetMulti([]string) (map[string]*Item, error)                      { return nil, c.err }
func (c errConn) Touch(string, int32) error                                        { return c.err }
func (c errConn) Delete(string) error                                              { return c.err }
func (c errConn) Increment(string, uint64) (uint64, error)                         { return 0, c.err }
func (c errConn) Decrement(string, uint64) (uint64, error)                         { return 0, c.err }
func (c errConn) Scan(*Item, interface{}) error                                    { return c.err }
func (c errConn) AddContext(context.Context, *Item) error                          { return c.err }
func (c errConn) SetContext(context.Context, *Item) error                          { return c.err }
func (c errConn) ReplaceContext(context.Context, *Item) error                      { return c.err }
func (c errConn) GetContext(context.Context, string) (*Item, error)                { return nil, c.err }
func (c errConn) DecrementContext(context.Context, string, uint64) (uint64, error) { return 0, c.err }
func (c errConn) CompareAndSwapContext(context.Context, *Item) error               { return c.err }
func (c errConn) TouchContext(context.Context, string, int32) error                { return c.err }
func (c errConn) DeleteContext(context.Context, string) error                      { return c.err }
func (c errConn) IncrementContext(context.Context, string, uint64) (uint64, error) { return 0, c.err }
func (c errConn) GetMultiContext(context.Context, []string) (map[string]*Item, error) {
	return nil, c.err
}

// RawItem item with FlagRAW flag.
//
// Expiration is the cache expiration time, in seconds: either a relative
// time from now (up to 1 month), or an absolute Unix epoch time.
// Zero means the Item has no expiration time.
func RawItem(key string, data []byte, flags uint32, expiration int32) *Item {
	return &Item{Key: key, Flags: flags | FlagRAW, Value: data, Expiration: expiration}
}

// JSONItem item with FlagJSON flag.
//
// Expiration is the cache expiration time, in seconds: either a relative
// time from now (up to 1 month), or an absolute Unix epoch time.
// Zero means the Item has no expiration time.
func JSONItem(key string, v interface{}, flags uint32, expiration int32) *Item {
	return &Item{Key: key, Flags: flags | FlagJSON, Object: v, Expiration: expiration}
}

// ProtobufItem item with FlagProtobuf flag.
//
// Expiration is the cache expiration time, in seconds: either a relative
// time from now (up to 1 month), or an absolute Unix epoch time.
// Zero means the Item has no expiration time.
func ProtobufItem(key string, message proto.Message, flags uint32, expiration int32) *Item {
	return &Item{Key: key, Flags: flags | FlagProtobuf, Object: message, Expiration: expiration}
}

func shrinkDeadline(ctx context.Context, timeout time.Duration) time.Time {
	timeoutTime := time.Now().Add(timeout)
	if deadline, ok := ctx.Deadline(); ok && timeoutTime.After(deadline) {
		return deadline
	}
	return timeoutTime
}
