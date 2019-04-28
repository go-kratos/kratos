package cache

import (
	"errors"
	"runtime"
	"sync"

	"go-common/library/log"
	"go-common/library/stat/prom"
)

var (
	// ErrFull cache internal chan full.
	ErrFull = errors.New("cache chan full")
	stats   = prom.BusinessInfoCount
)

// Cache async save data by chan.
type Cache struct {
	ch     chan func()
	worker int
	waiter sync.WaitGroup
}

// Deprecated: use library/sync/pipeline/fanout instead.
func New(worker, size int) *Cache {
	if worker <= 0 {
		worker = 1
	}
	c := &Cache{
		ch:     make(chan func(), size),
		worker: worker,
	}
	c.waiter.Add(worker)
	for i := 0; i < worker; i++ {
		go c.proc()
	}
	return c
}

func (c *Cache) proc() {
	defer c.waiter.Done()
	for {
		f := <-c.ch
		if f == nil {
			return
		}
		wrapFunc(f)()
		stats.State("cache_channel", int64(len(c.ch)))
	}
}

func wrapFunc(f func()) (res func()) {
	res = func() {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64*1024)
				buf = buf[:runtime.Stack(buf, false)]
				log.Error("panic in cache proc, err: %s, stack: %s", r, buf)
			}
		}()
		f()
	}
	return
}

// Save save a callback cache func.
func (c *Cache) Save(f func()) (err error) {
	if f == nil {
		return
	}
	select {
	case c.ch <- f:
	default:
		err = ErrFull
	}
	stats.State("cache_channel", int64(len(c.ch)))
	return
}

// Close close cache
func (c *Cache) Close() (err error) {
	for i := 0; i < c.worker; i++ {
		c.ch <- nil
	}
	c.waiter.Wait()
	return
}
