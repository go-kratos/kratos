package pool

import (
	"context"
	"io"
	"sync"
	"time"
)

var _ Pool = &Slice{}

// Slice .
type Slice struct {
	// New is an application supplied function for creating and configuring a
	// item.
	//
	// The item returned from new must not be in a special state
	// (subscribed to pubsub channel, transaction started, ...).
	New  func(ctx context.Context) (io.Closer, error)
	stop func() // stop cancels the item opener.

	// mu protects fields defined below.
	mu           sync.Mutex
	freeItem     []*item
	itemRequests map[uint64]chan item
	nextRequest  uint64 // Next key to use in itemRequests.
	active       int    // number of opened and pending open items
	// Used to signal the need for new items
	// a goroutine running itemOpener() reads on this chan and
	// maybeOpenNewItems sends on the chan (one send per needed item)
	// It is closed during db.Close(). The close tells the itemOpener
	// goroutine to exit.
	openerCh  chan struct{}
	closed    bool
	cleanerCh chan struct{}

	// Config pool configuration
	conf *Config
}

// NewSlice creates a new pool.
func NewSlice(c *Config) *Slice {
	// check Config
	if c == nil || c.Active < c.Idle {
		panic("config nil or Idle Must <= Active")
	}
	ctx, cancel := context.WithCancel(context.Background())
	// new pool
	p := &Slice{
		conf:         c,
		stop:         cancel,
		itemRequests: make(map[uint64]chan item),
		openerCh:     make(chan struct{}, 1000000),
	}
	p.startCleanerLocked(time.Duration(c.IdleTimeout))

	go p.itemOpener(ctx)
	return p
}

// Reload reload config.
func (p *Slice) Reload(c *Config) error {
	p.mu.Lock()
	p.startCleanerLocked(time.Duration(c.IdleTimeout))
	p.setActive(c.Active)
	p.setIdle(c.Idle)
	p.conf = c
	p.mu.Unlock()
	return nil
}

// Get returns a newly-opened or cached *item.
func (p *Slice) Get(ctx context.Context) (io.Closer, error) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil, ErrPoolClosed
	}
	idleTimeout := time.Duration(p.conf.IdleTimeout)
	// Prefer a free item, if possible.
	numFree := len(p.freeItem)
	for numFree > 0 {
		i := p.freeItem[0]
		copy(p.freeItem, p.freeItem[1:])
		p.freeItem = p.freeItem[:numFree-1]
		p.mu.Unlock()
		if i.expired(idleTimeout) {
			i.close()
			p.mu.Lock()
			p.release()
		} else {
			return i.c, nil
		}
		numFree = len(p.freeItem)
	}

	// Out of free items or we were asked not to use one. If we're not
	// allowed to open any more items, make a request and wait.
	if p.conf.Active > 0 && p.active >= p.conf.Active {
		// check WaitTimeout and return directly
		if p.conf.WaitTimeout == 0 && !p.conf.Wait {
			p.mu.Unlock()
			return nil, ErrPoolExhausted
		}
		// Make the item channel. It's buffered so that the
		// itemOpener doesn't block while waiting for the req to be read.
		req := make(chan item, 1)
		reqKey := p.nextRequestKeyLocked()
		p.itemRequests[reqKey] = req
		wt := p.conf.WaitTimeout
		p.mu.Unlock()

		// reset context timeout
		if wt > 0 {
			var cancel func()
			_, ctx, cancel = wt.Shrink(ctx)
			defer cancel()
		}
		// Timeout the item request with the context.
		select {
		case <-ctx.Done():
			// Remove the item request and ensure no value has been sent
			// on it after removing.
			p.mu.Lock()
			delete(p.itemRequests, reqKey)
			p.mu.Unlock()
			return nil, ctx.Err()
		case ret, ok := <-req:
			if !ok {
				return nil, ErrPoolClosed
			}
			if ret.expired(idleTimeout) {
				ret.close()
				p.mu.Lock()
				p.release()
			} else {
				return ret.c, nil
			}
		}
	}

	p.active++ // optimistically
	p.mu.Unlock()
	c, err := p.New(ctx)
	if err != nil {
		p.mu.Lock()
		p.release()
		p.mu.Unlock()
		return nil, err
	}
	return c, nil
}

// Put adds a item to the p's free pool.
// err is optionally the last error that occurred on this item.
func (p *Slice) Put(ctx context.Context, c io.Closer, forceClose bool) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if forceClose {
		p.release()
		return c.Close()
	}
	added := p.putItemLocked(c)
	if !added {
		p.active--
		return c.Close()
	}
	return nil
}

// Satisfy a item or put the item in the idle pool and return true
// or return false.
// putItemLocked will satisfy a item if there is one, or it will
// return the *item to the freeItem list if err == nil and the idle
// item limit will not be exceeded.
// If err != nil, the value of i is ignored.
// If err == nil, then i must not equal nil.
// If a item was fulfilled or the *item was placed in the
// freeItem list, then true is returned, otherwise false is returned.
func (p *Slice) putItemLocked(c io.Closer) bool {
	if p.closed {
		return false
	}
	if p.conf.Active > 0 && p.active > p.conf.Active {
		return false
	}
	i := item{
		c:         c,
		createdAt: nowFunc(),
	}
	if l := len(p.itemRequests); l > 0 {
		var req chan item
		var reqKey uint64
		for reqKey, req = range p.itemRequests {
			break
		}
		delete(p.itemRequests, reqKey) // Remove from pending requests.
		req <- i
		return true
	} else if !p.closed && p.maxIdleItemsLocked() > len(p.freeItem) {
		p.freeItem = append(p.freeItem, &i)
		return true
	}
	return false
}

// Runs in a separate goroutine, opens new item when requested.
func (p *Slice) itemOpener(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.openerCh:
			p.openNewItem(ctx)
		}
	}
}

func (p *Slice) maybeOpenNewItems() {
	numRequests := len(p.itemRequests)
	if p.conf.Active > 0 {
		numCanOpen := p.conf.Active - p.active
		if numRequests > numCanOpen {
			numRequests = numCanOpen
		}
	}
	for numRequests > 0 {
		p.active++ // optimistically
		numRequests--
		if p.closed {
			return
		}
		p.openerCh <- struct{}{}
	}
}

// openNewItem one new item
func (p *Slice) openNewItem(ctx context.Context) {
	// maybeOpenNewConnctions has already executed p.active++ before it sent
	// on p.openerCh. This function must execute p.active-- if the
	// item fails or is closed before returning.
	c, err := p.New(ctx)
	p.mu.Lock()
	defer p.mu.Unlock()
	if err != nil {
		p.release()
		return
	}
	if !p.putItemLocked(c) {
		p.active--
		c.Close()
	}
}

// setIdle sets the maximum number of items in the idle
// item pool.
//
// If MaxOpenConns is greater than 0 but less than the new IdleConns
// then the new IdleConns will be reduced to match the MaxOpenConns limit
//
// If n <= 0, no idle items are retained.
func (p *Slice) setIdle(n int) {
	p.mu.Lock()
	if n > 0 {
		p.conf.Idle = n
	} else {
		// No idle items.
		p.conf.Idle = -1
	}
	// Make sure maxIdle doesn't exceed maxOpen
	if p.conf.Active > 0 && p.maxIdleItemsLocked() > p.conf.Active {
		p.conf.Idle = p.conf.Active
	}
	var closing []*item
	idleCount := len(p.freeItem)
	maxIdle := p.maxIdleItemsLocked()
	if idleCount > maxIdle {
		closing = p.freeItem[maxIdle:]
		p.freeItem = p.freeItem[:maxIdle]
	}
	p.mu.Unlock()
	for _, c := range closing {
		c.close()
	}
}

// setActive sets the maximum number of open items to the database.
//
// If IdleConns is greater than 0 and the new MaxOpenConns is less than
// IdleConns, then IdleConns will be reduced to match the new
// MaxOpenConns limit
//
// If n <= 0, then there is no limit on the number of open items.
// The default is 0 (unlimited).
func (p *Slice) setActive(n int) {
	p.mu.Lock()
	p.conf.Active = n
	if n < 0 {
		p.conf.Active = 0
	}
	syncIdle := p.conf.Active > 0 && p.maxIdleItemsLocked() > p.conf.Active
	p.mu.Unlock()
	if syncIdle {
		p.setIdle(n)
	}
}

// startCleanerLocked starts itemCleaner if needed.
func (p *Slice) startCleanerLocked(d time.Duration) {
	if d <= 0 {
		// if set 0, staleCleaner() will return directly
		return
	}
	if d < time.Duration(p.conf.IdleTimeout) && p.cleanerCh != nil {
		select {
		case p.cleanerCh <- struct{}{}:
		default:
		}
	}
	// run only one, clean stale items.
	if p.cleanerCh == nil {
		p.cleanerCh = make(chan struct{}, 1)
		go p.staleCleaner(time.Duration(p.conf.IdleTimeout))
	}
}

func (p *Slice) staleCleaner(d time.Duration) {
	const minInterval = 100 * time.Millisecond

	if d < minInterval {
		d = minInterval
	}
	t := time.NewTimer(d)

	for {
		select {
		case <-t.C:
		case <-p.cleanerCh: // maxLifetime was changed or db was closed.
		}
		p.mu.Lock()
		d = time.Duration(p.conf.IdleTimeout)
		if p.closed || d <= 0 {
			p.mu.Unlock()
			return
		}

		expiredSince := nowFunc().Add(-d)
		var closing []*item
		for i := 0; i < len(p.freeItem); i++ {
			c := p.freeItem[i]
			if c.createdAt.Before(expiredSince) {
				closing = append(closing, c)
				p.active--
				last := len(p.freeItem) - 1
				p.freeItem[i] = p.freeItem[last]
				p.freeItem[last] = nil
				p.freeItem = p.freeItem[:last]
				i--
			}
		}
		p.mu.Unlock()

		for _, c := range closing {
			c.close()
		}

		if d < minInterval {
			d = minInterval
		}
		t.Reset(d)
	}
}

// nextRequestKeyLocked returns the next item request key.
// It is assumed that nextRequest will not overflow.
func (p *Slice) nextRequestKeyLocked() uint64 {
	next := p.nextRequest
	p.nextRequest++
	return next
}

const defaultIdleItems = 2

func (p *Slice) maxIdleItemsLocked() int {
	n := p.conf.Idle
	switch {
	case n == 0:
		return defaultIdleItems
	case n < 0:
		return 0
	default:
		return n
	}
}

func (p *Slice) release() {
	p.active--
	p.maybeOpenNewItems()
}

// Close close pool.
func (p *Slice) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	if p.cleanerCh != nil {
		close(p.cleanerCh)
	}
	var err error
	for _, i := range p.freeItem {
		i.close()
	}
	p.freeItem = nil
	p.closed = true
	for _, req := range p.itemRequests {
		close(req)
	}
	p.mu.Unlock()
	p.stop()
	return err
}
