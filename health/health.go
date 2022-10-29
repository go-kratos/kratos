package health

import (
	"context"
	"sync"
)

type CheckerMgr struct {
	checkers  map[string]*checker
	ctx       context.Context
	cancel    func()
	watchers  map[uint64]chan string
	watcherID uint64

	lock sync.RWMutex
}

func New(ctx context.Context) *CheckerMgr {
	c, cancel := context.WithCancel(ctx)
	return &CheckerMgr{
		checkers: make(map[string]*checker),
		ctx:      c,
		cancel:   cancel,
		lock:     sync.RWMutex{},
		watchers: map[uint64]chan string{},
	}
}

func (c *CheckerMgr) Start() {
	for _, v := range c.checkers {
		cv := v
		go func() {
			cv.run(c.ctx)
		}()
	}
}

func (c *CheckerMgr) Stop() {
	c.cancel()
	c.closeAllWatcher()
}

type StatusResult struct {
	Name string
	CheckerStatus
}

// GetStatus
//
//	if name is nil return all status
func (c *CheckerMgr) GetStatus(name ...string) []StatusResult {
	status := make([]StatusResult, 0, len(name))
	if len(name) == 0 {
		for _, v := range c.checkers {
			status = append(status, StatusResult{
				Name:          v.Name,
				CheckerStatus: v.getStatus(),
			})
		}
	} else {
		for _, n := range name {
			if v, ok := c.checkers[n]; ok {
				status = append(status, StatusResult{
					Name:          v.Name,
					CheckerStatus: v.getStatus(),
				})
			}
		}
	}
	return status
}

func (c *CheckerMgr) RegisterChecker(checker2 *checker) {
	c.checkers[checker2.Name] = checker2
	checker2.setWatcher(c)
}

func (c *CheckerMgr) Watch(name string) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	for _, ch := range c.watchers {
		select {
		case ch <- name:
		default:
		}
	}
}

type WatcherResult struct {
	id uint64
	Ch <-chan string
	c  *CheckerMgr
}

func (w *WatcherResult) Close() {
	w.c.closeWatcher(w.id)
}

func (c *CheckerMgr) NewWatcher() WatcherResult {
	c.lock.Lock()
	wID := c.watcherID
	c.watcherID++
	ch := make(chan string, 1)
	c.watchers[wID] = ch
	c.lock.Unlock()
	return WatcherResult{
		id: wID,
		Ch: ch,
		c:  c,
	}
}

func (c *CheckerMgr) closeWatcher(wID uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	ch, ok := c.watchers[wID]
	if !ok {
		return
	}
	close(ch)
	delete(c.watchers, wID)
}

func (c *CheckerMgr) closeAllWatcher() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for k, v := range c.watchers {
		close(v)
		delete(c.watchers, k)
	}
}
