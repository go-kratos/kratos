package health

import (
	"context"
	"sync"
)

type CheckerMgr struct {
	checkers  map[string]CheckerHandler
	ctx       context.Context
	cancel    func()
	watchers  map[uint64]chan string
	watcherID uint64

	lock sync.RWMutex
}

func New(ctx context.Context) *CheckerMgr {
	c, cancel := context.WithCancel(ctx)
	return &CheckerMgr{
		checkers: make(map[string]CheckerHandler),
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
				Name:          v.getName(),
				CheckerStatus: v.getStatus(),
			})
		}
	} else {
		for _, n := range name {
			if v, ok := c.checkers[n]; ok {
				status = append(status, StatusResult{
					Name:          v.getName(),
					CheckerStatus: v.getStatus(),
				})
			}
		}
	}
	return status
}

func (c *CheckerMgr) RegisterChecker(checker CheckerHandler) {
	c.checkers[checker.getName()] = checker
	checker.setNotifier(c)
}

func (c *CheckerMgr) notify(name string) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	for _, ch := range c.watchers {
		select {
		case ch <- name:
		default:
		}
	}
}

type Watcher struct {
	id uint64
	Ch <-chan string
	c  *CheckerMgr
}

func (w *Watcher) Close() {
	w.c.closeWatcher(w.id)
}

func (c *CheckerMgr) NewWatcher() Watcher {
	c.lock.Lock()
	wID := c.watcherID
	c.watcherID++
	ch := make(chan string, 1)
	c.watchers[wID] = ch
	c.lock.Unlock()
	return Watcher{
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
