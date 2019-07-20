package memcache

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/bilibili/kratos/pkg/container/pool"
)

// Pool memcache connection pool struct.
// Deprecated: Use Memcache instead
type Pool struct {
	p pool.Pool
	c *Config
}

// NewPool new a memcache conn pool.
// Deprecated: Use New instead
func NewPool(cfg *Config) (p *Pool) {
	if cfg.DialTimeout <= 0 || cfg.ReadTimeout <= 0 || cfg.WriteTimeout <= 0 {
		panic("must config memcache timeout")
	}
	p1 := pool.NewList(cfg.Config)
	cnop := DialConnectTimeout(time.Duration(cfg.DialTimeout))
	rdop := DialReadTimeout(time.Duration(cfg.ReadTimeout))
	wrop := DialWriteTimeout(time.Duration(cfg.WriteTimeout))
	p1.New = func(ctx context.Context) (io.Closer, error) {
		conn, err := Dial(cfg.Proto, cfg.Addr, cnop, rdop, wrop)
		return newTraceConn(conn, fmt.Sprintf("%s://%s", cfg.Proto, cfg.Addr)), err
	}
	p = &Pool{p: p1, c: cfg}
	return
}

// Get gets a connection. The application must close the returned connection.
// This method always returns a valid connection so that applications can defer
// error handling to the first use of the connection. If there is an error
// getting an underlying connection, then the connection Err, Do, Send, Flush
// and Receive methods return that error.
func (p *Pool) Get(ctx context.Context) Conn {
	c, err := p.p.Get(ctx)
	if err != nil {
		return errConn{err}
	}
	c1, _ := c.(Conn)
	return &poolConn{p: p, c: c1, ctx: ctx}
}

// Close release the resources used by the pool.
func (p *Pool) Close() error {
	return p.p.Close()
}

type poolConn struct {
	c   Conn
	p   *Pool
	ctx context.Context
}

func (pc *poolConn) pstat(key string, t time.Time, err error) {
	_metricReqDur.Observe(int64(time.Since(t)/time.Millisecond), pc.p.c.Name, pc.p.c.Addr, key)
	if err != nil {
		if msg := pc.formatErr(err); msg != "" {
			_metricReqErr.Inc(pc.p.c.Name, pc.p.c.Addr, key, msg)
		}
		return
	}
	_metricHits.Inc(pc.p.c.Name, pc.p.c.Addr)
}

func (pc *poolConn) Close() error {
	c := pc.c
	if _, ok := c.(errConn); ok {
		return nil
	}
	pc.c = errConn{ErrConnClosed}
	pc.p.p.Put(context.Background(), c, c.Err() != nil)
	return nil
}

func (pc *poolConn) Err() error {
	return pc.c.Err()
}

func (pc *poolConn) Set(item *Item) (err error) {
	return pc.SetContext(pc.ctx, item)
}

func (pc *poolConn) Add(item *Item) (err error) {
	return pc.AddContext(pc.ctx, item)
}

func (pc *poolConn) Replace(item *Item) (err error) {
	return pc.ReplaceContext(pc.ctx, item)
}

func (pc *poolConn) CompareAndSwap(item *Item) (err error) {
	return pc.CompareAndSwapContext(pc.ctx, item)
}

func (pc *poolConn) Get(key string) (r *Item, err error) {
	return pc.GetContext(pc.ctx, key)
}

func (pc *poolConn) GetMulti(keys []string) (res map[string]*Item, err error) {
	return pc.GetMultiContext(pc.ctx, keys)
}

func (pc *poolConn) Touch(key string, timeout int32) (err error) {
	return pc.TouchContext(pc.ctx, key, timeout)
}

func (pc *poolConn) Scan(item *Item, v interface{}) error {
	return pc.c.Scan(item, v)
}

func (pc *poolConn) Delete(key string) (err error) {
	return pc.DeleteContext(pc.ctx, key)
}

func (pc *poolConn) Increment(key string, delta uint64) (newValue uint64, err error) {
	return pc.IncrementContext(pc.ctx, key, delta)
}

func (pc *poolConn) Decrement(key string, delta uint64) (newValue uint64, err error) {
	return pc.DecrementContext(pc.ctx, key, delta)
}

func (pc *poolConn) AddContext(ctx context.Context, item *Item) error {
	now := time.Now()
	err := pc.c.AddContext(ctx, item)
	pc.pstat("add", now, err)
	return err
}

func (pc *poolConn) SetContext(ctx context.Context, item *Item) error {
	now := time.Now()
	err := pc.c.SetContext(ctx, item)
	pc.pstat("set", now, err)
	return err
}

func (pc *poolConn) ReplaceContext(ctx context.Context, item *Item) error {
	now := time.Now()
	err := pc.c.ReplaceContext(ctx, item)
	pc.pstat("replace", now, err)
	return err
}

func (pc *poolConn) GetContext(ctx context.Context, key string) (*Item, error) {
	now := time.Now()
	item, err := pc.c.Get(key)
	pc.pstat("get", now, err)
	return item, err
}

func (pc *poolConn) GetMultiContext(ctx context.Context, keys []string) (map[string]*Item, error) {
	// if keys is empty slice returns empty map direct
	if len(keys) == 0 {
		return make(map[string]*Item), nil
	}
	now := time.Now()
	items, err := pc.c.GetMulti(keys)
	pc.pstat("gets", now, err)
	return items, err
}

func (pc *poolConn) DeleteContext(ctx context.Context, key string) error {
	now := time.Now()
	err := pc.c.Delete(key)
	pc.pstat("delete", now, err)
	return err
}

func (pc *poolConn) IncrementContext(ctx context.Context, key string, delta uint64) (uint64, error) {
	now := time.Now()
	newValue, err := pc.c.IncrementContext(ctx, key, delta)
	pc.pstat("increment", now, err)
	return newValue, err
}

func (pc *poolConn) DecrementContext(ctx context.Context, key string, delta uint64) (uint64, error) {
	now := time.Now()
	newValue, err := pc.c.DecrementContext(ctx, key, delta)
	pc.pstat("decrement", now, err)
	return newValue, err
}

func (pc *poolConn) CompareAndSwapContext(ctx context.Context, item *Item) error {
	now := time.Now()
	err := pc.c.CompareAndSwap(item)
	pc.pstat("cas", now, err)
	return err
}

func (pc *poolConn) TouchContext(ctx context.Context, key string, seconds int32) error {
	now := time.Now()
	err := pc.c.Touch(key, seconds)
	pc.pstat("touch", now, err)
	return err
}
