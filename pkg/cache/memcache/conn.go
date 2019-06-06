package memcache

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	pkgerr "github.com/pkg/errors"
)

const (
	// 1024*1024 - 1, set error???
	_largeValue = 1000 * 1000 // 1MB
)

// low level connection that implement memcache protocol provide basic operation.
type protocolConn interface {
	Populate(ctx context.Context, cmd string, key string, flags uint32, expiration int32, cas uint64, data []byte) error
	Get(ctx context.Context, key string) (*Item, error)
	GetMulti(ctx context.Context, keys ...string) (map[string]*Item, error)
	Touch(ctx context.Context, key string, expire int32) error
	IncrDecr(ctx context.Context, cmd, key string, delta uint64) (uint64, error)
	Delete(ctx context.Context, key string) error
	Close() error
	Err() error
}

// DialOption specifies an option for dialing a Memcache server.
type DialOption struct {
	f func(*dialOptions)
}

type dialOptions struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	protocol     string
	dial         func(network, addr string) (net.Conn, error)
}

// DialReadTimeout specifies the timeout for reading a single command reply.
func DialReadTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.readTimeout = d
	}}
}

// DialWriteTimeout specifies the timeout for writing a single command.
func DialWriteTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		do.writeTimeout = d
	}}
}

// DialConnectTimeout specifies the timeout for connecting to the Memcache server.
func DialConnectTimeout(d time.Duration) DialOption {
	return DialOption{func(do *dialOptions) {
		dialer := net.Dialer{Timeout: d}
		do.dial = dialer.Dial
	}}
}

// DialNetDial specifies a custom dial function for creating TCP
// connections. If this option is left out, then net.Dial is
// used. DialNetDial overrides DialConnectTimeout.
func DialNetDial(dial func(network, addr string) (net.Conn, error)) DialOption {
	return DialOption{func(do *dialOptions) {
		do.dial = dial
	}}
}

// Dial connects to the Memcache server at the given network and
// address using the specified options.
func Dial(network, address string, options ...DialOption) (Conn, error) {
	do := dialOptions{
		dial: net.Dial,
	}
	for _, option := range options {
		option.f(&do)
	}
	netConn, err := do.dial(network, address)
	if err != nil {
		return nil, pkgerr.WithStack(err)
	}
	pconn, err := newASCIIConn(netConn, do.readTimeout, do.writeTimeout)
	return &conn{pconn: pconn, ed: newEncodeDecoder()}, nil
}

type conn struct {
	// low level connection.
	pconn protocolConn
	ed    *encodeDecode
}

func (c *conn) Close() error {
	return c.pconn.Close()
}

func (c *conn) Err() error {
	return c.pconn.Err()
}

func (c *conn) AddContext(ctx context.Context, item *Item) error {
	return c.populate(ctx, "add", item)
}

func (c *conn) SetContext(ctx context.Context, item *Item) error {
	return c.populate(ctx, "set", item)
}

func (c *conn) ReplaceContext(ctx context.Context, item *Item) error {
	return c.populate(ctx, "replace", item)
}

func (c *conn) CompareAndSwapContext(ctx context.Context, item *Item) error {
	return c.populate(ctx, "cas", item)
}

func (c *conn) populate(ctx context.Context, cmd string, item *Item) error {
	if !legalKey(item.Key) {
		return ErrMalformedKey
	}
	data, err := c.ed.encode(item)
	if err != nil {
		return err
	}
	length := len(data)
	if length < _largeValue {
		return c.pconn.Populate(ctx, cmd, item.Key, item.Flags, item.Expiration, item.cas, data)
	}
	count := length/_largeValue + 1
	if err = c.pconn.Populate(ctx, cmd, item.Key, item.Flags|flagLargeValue, item.Expiration, item.cas, []byte(strconv.Itoa(length))); err != nil {
		return err
	}
	var chunk []byte
	for i := 1; i <= count; i++ {
		if i == count {
			chunk = data[_largeValue*(count-1):]
		} else {
			chunk = data[_largeValue*(i-1) : _largeValue*i]
		}
		key := fmt.Sprintf("%s%d", item.Key, i)
		if err = c.pconn.Populate(ctx, cmd, key, item.Flags, item.Expiration, item.cas, chunk); err != nil {
			return err
		}
	}
	return nil
}

func (c *conn) GetContext(ctx context.Context, key string) (*Item, error) {
	if !legalKey(key) {
		return nil, ErrMalformedKey
	}
	result, err := c.pconn.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if result.Flags&flagLargeValue != flagLargeValue {
		return result, err
	}
	return c.getLargeItem(ctx, result)
}

func (c *conn) getLargeItem(ctx context.Context, result *Item) (*Item, error) {
	length, err := strconv.Atoi(string(result.Value))
	if err != nil {
		return nil, err
	}
	count := length/_largeValue + 1
	keys := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		keys = append(keys, fmt.Sprintf("%s%d", result.Key, i))
	}
	var results map[string]*Item
	if results, err = c.pconn.GetMulti(ctx, keys...); err != nil {
		return nil, err
	}
	if len(results) < count {
		return nil, ErrNotFound
	}
	result.Value = make([]byte, 0, length)
	for _, k := range keys {
		ti := results[k]
		if ti == nil || ti.Value == nil {
			return nil, ErrNotFound
		}
		result.Value = append(result.Value, ti.Value...)
	}
	result.Flags = result.Flags ^ flagLargeValue
	return result, nil
}

func (c *conn) GetMultiContext(ctx context.Context, keys []string) (map[string]*Item, error) {
	// TODO: move to protocolConn?
	for _, key := range keys {
		if !legalKey(key) {
			return nil, ErrMalformedKey
		}
	}
	results, err := c.pconn.GetMulti(ctx, keys...)
	if err != nil {
		return results, err
	}
	for k, v := range results {
		if v.Flags&flagLargeValue != flagLargeValue {
			continue
		}
		if v, err = c.getLargeItem(ctx, v); err != nil {
			return results, err
		}
		results[k] = v
	}
	return results, nil
}

func (c *conn) DeleteContext(ctx context.Context, key string) error {
	if !legalKey(key) {
		return ErrMalformedKey
	}
	return c.pconn.Delete(ctx, key)
}

func (c *conn) IncrementContext(ctx context.Context, key string, delta uint64) (uint64, error) {
	if !legalKey(key) {
		return 0, ErrMalformedKey
	}
	return c.pconn.IncrDecr(ctx, "incr", key, delta)
}

func (c *conn) DecrementContext(ctx context.Context, key string, delta uint64) (uint64, error) {
	if !legalKey(key) {
		return 0, ErrMalformedKey
	}
	return c.pconn.IncrDecr(ctx, "decr", key, delta)
}

func (c *conn) TouchContext(ctx context.Context, key string, seconds int32) error {
	if !legalKey(key) {
		return ErrMalformedKey
	}
	return c.pconn.Touch(ctx, key, seconds)
}

func (c *conn) Add(item *Item) error {
	return c.AddContext(context.TODO(), item)
}

func (c *conn) Set(item *Item) error {
	return c.SetContext(context.TODO(), item)
}

func (c *conn) Replace(item *Item) error {
	return c.ReplaceContext(context.TODO(), item)
}

func (c *conn) Get(key string) (*Item, error) {
	return c.GetContext(context.TODO(), key)
}

func (c *conn) GetMulti(keys []string) (map[string]*Item, error) {
	return c.GetMultiContext(context.TODO(), keys)
}

func (c *conn) Delete(key string) error {
	return c.DeleteContext(context.TODO(), key)
}

func (c *conn) Increment(key string, delta uint64) (newValue uint64, err error) {
	return c.IncrementContext(context.TODO(), key, delta)
}

func (c *conn) Decrement(key string, delta uint64) (newValue uint64, err error) {
	return c.DecrementContext(context.TODO(), key, delta)
}

func (c *conn) CompareAndSwap(item *Item) error {
	return c.CompareAndSwapContext(context.TODO(), item)
}

func (c *conn) Touch(key string, seconds int32) (err error) {
	return c.TouchContext(context.TODO(), key, seconds)
}

func (c *conn) Scan(item *Item, v interface{}) (err error) {
	return pkgerr.WithStack(c.ed.decode(item, v))
}
