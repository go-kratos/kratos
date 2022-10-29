package health

import (
	"context"
	"sync"
	"time"
)

type Checker interface {
	Check(ctx context.Context) (interface{}, error)
}

type Watcher interface {
	Watch(string)
}

type checker struct {
	Name         string
	intervalTime time.Duration
	timeout      time.Duration
	Checker
	CheckerStatus
	*sync.RWMutex
	Watcher
}

func NewChecker(name string, ch Checker, interval, timeout time.Duration) *checker {
	return &checker{
		Name:          name,
		intervalTime:  interval,
		timeout:       timeout,
		Checker:       ch,
		CheckerStatus: CheckerStatus{},
		RWMutex:       &sync.RWMutex{},
		Watcher:       nil,
	}
}

func (c *checker) setWatcher(w Watcher) {
	c.Watcher = w
}

func (c *checker) check(ctx context.Context) bool {
	defer func() {
		recover()
	}()

	var cancel func()
	if c.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	detail, err := c.Check(ctx)
	status := StatusUp
	if err != nil {
		status = StatusDown
	}

	c.Lock()
	defer c.Unlock()
	old := c.CheckerStatus
	c.CheckerStatus = CheckerStatus{
		Status: status,
		Detail: detail,
		Err:    err,
	}
	if c.CheckerStatus == old {
		return false
	}
	return true
}

func (c *checker) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if c.check(ctx) {
			//notify
			if c.Watcher != nil {
				c.Watcher.Watch(c.Name)
			}
		}
		time.Sleep(c.intervalTime)
	}
}

func (c *checker) getStatus() CheckerStatus {
	c.RLock()
	defer c.RUnlock()
	return c.CheckerStatus
}
