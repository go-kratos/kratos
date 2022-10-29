package health

import (
	"context"
	"sync"
	"time"
)

type Checker interface {
	Check(ctx context.Context) (interface{}, error)
}

type Notifier interface {
	notify(string)
}

type CheckerHandler interface {
	setNotifier(w Notifier)
	run(ctx context.Context)
	getStatus() CheckerStatus
	getName() string
}

type checkerHandler struct {
	Name         string
	intervalTime time.Duration
	timeout      time.Duration
	Checker
	CheckerStatus
	*sync.RWMutex
	Notifier
}

func NewChecker(name string, ch Checker, interval, timeout time.Duration) CheckerHandler {
	return &checkerHandler{
		Name:          name,
		intervalTime:  interval,
		timeout:       timeout,
		Checker:       ch,
		CheckerStatus: CheckerStatus{},
		RWMutex:       &sync.RWMutex{},
		Notifier:      nil,
	}
}

func (c *checkerHandler) setNotifier(w Notifier) {
	c.Notifier = w
}

func (c *checkerHandler) getName() string {
	return c.Name
}

func (c *checkerHandler) check(ctx context.Context) bool {
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

func (c *checkerHandler) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if c.check(ctx) {
			//notify
			if c.Notifier != nil {
				c.Notifier.notify(c.Name)
			}
		}
		time.Sleep(c.intervalTime)
	}
}

func (c *checkerHandler) getStatus() CheckerStatus {
	c.RLock()
	defer c.RUnlock()
	return c.CheckerStatus
}
