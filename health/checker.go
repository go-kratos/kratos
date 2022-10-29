package health

import (
	"context"
	"sync"
	"time"
)

type Checker interface {
	Check(ctx context.Context) (interface{}, error)
}

type Watcher func(string)

type checker struct {
	Name         string
	intervalTime time.Duration
	timeout      time.Duration
	Checker
	CheckerStatus
	sync.RWMutex
	Watcher
}

func NewChecker(name string, checker Checker) {

}

func (c *checker) check(ctx context.Context) bool {
	defer func() {
		recover()
	}()
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
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
			//发送改变通知
			c.Watcher(c.Name)
		}
		time.Sleep(c.intervalTime)
	}
}

func (c *checker) getStatus() CheckerStatus {
	c.RLock()
	defer c.RUnlock()
	return c.CheckerStatus
}
