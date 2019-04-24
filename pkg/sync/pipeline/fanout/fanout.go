package fanout

import (
	"context"
	"errors"
	"runtime"
	"sync"

	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/net/trace"
	"go-common/library/stat/prom"
)

var (
	// ErrFull chan full.
	ErrFull   = errors.New("fanout: chan full")
	stats     = prom.BusinessInfoCount
	traceTags = []trace.Tag{
		trace.Tag{Key: trace.TagSpanKind, Value: "background"},
		trace.Tag{Key: trace.TagComponent, Value: "sync/pipeline/fanout"},
	}
)

type options struct {
	worker int
	buffer int
}

// Option fanout option
type Option func(*options)

// Worker specifies the worker of fanout
func Worker(n int) Option {
	if n <= 0 {
		panic("fanout: worker should > 0")
	}
	return func(o *options) {
		o.worker = n
	}
}

// Buffer specifies the buffer of fanout
func Buffer(n int) Option {
	if n <= 0 {
		panic("fanout: buffer should > 0")
	}
	return func(o *options) {
		o.buffer = n
	}
}

type item struct {
	f   func(c context.Context)
	ctx context.Context
}

// Fanout async consume data from chan.
type Fanout struct {
	name    string
	ch      chan item
	options *options
	waiter  sync.WaitGroup

	ctx    context.Context
	cancel func()
}

// New new a fanout struct.
func New(name string, opts ...Option) *Fanout {
	if name == "" {
		name = "fanout"
	}
	o := &options{
		worker: 1,
		buffer: 1024,
	}
	for _, op := range opts {
		op(o)
	}
	c := &Fanout{
		ch:      make(chan item, o.buffer),
		name:    name,
		options: o,
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.waiter.Add(o.worker)
	for i := 0; i < o.worker; i++ {
		go c.proc()
	}
	return c
}

func (c *Fanout) proc() {
	defer c.waiter.Done()
	for {
		select {
		case t := <-c.ch:
			wrapFunc(t.f)(t.ctx)
			stats.State(c.name+"_channel", int64(len(c.ch)))
		case <-c.ctx.Done():
			return
		}
	}
}

func wrapFunc(f func(c context.Context)) (res func(context.Context)) {
	res = func(ctx context.Context) {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64*1024)
				buf = buf[:runtime.Stack(buf, false)]
				log.Error("panic in fanout proc, err: %s, stack: %s", r, buf)
			}
		}()
		f(ctx)
		if tr, ok := trace.FromContext(ctx); ok {
			tr.Finish(nil)
		}
	}
	return
}

// Do save a callback func.
func (c *Fanout) Do(ctx context.Context, f func(ctx context.Context)) (err error) {
	if f == nil || c.ctx.Err() != nil {
		return c.ctx.Err()
	}
	nakeCtx := metadata.WithContext(ctx)
	if tr, ok := trace.FromContext(ctx); ok {
		tr = tr.Fork("", "Fanout:Do").SetTag(traceTags...)
		nakeCtx = trace.NewContext(nakeCtx, tr)
	}
	select {
	case c.ch <- item{f: f, ctx: nakeCtx}:
	default:
		err = ErrFull
	}
	stats.State(c.name+"_channel", int64(len(c.ch)))
	return
}

// Close close fanout
func (c *Fanout) Close() error {
	if err := c.ctx.Err(); err != nil {
		return err
	}
	c.cancel()
	c.waiter.Wait()
	return nil
}
