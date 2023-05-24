package health

import (
	"context"
	"fmt"
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

func NewChecker(name string, ch Checker, opt ...CheckOption) CheckerHandler {
	c := &checkerHandler{
		Name:          name,
		Checker:       ch,
		CheckerStatus: CheckerStatus{},
		RWMutex:       &sync.RWMutex{},
		Notifier:      nil,
	}
	for _, v := range opt {
		v(c)
	}
	return c
}

func (c *checkerHandler) setNotifier(w Notifier) {
	c.Notifier = w
}

func (c *checkerHandler) getName() string {
	return c.Name
}

func (c *checkerHandler) check(ctx context.Context) (r bool) {
	defer func() {
		if err := recover(); err != nil {
			r = true
			c.Lock()
			c.CheckerStatus.Err = fmt.Errorf("%v", err)
			c.Unlock()
		}
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

	return c.CheckerStatus != old
}

func (c *checkerHandler) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		if c.check(ctx) {
			// notify
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

type CheckOption func(handler *checkerHandler)

func WithInterval(interval time.Duration) CheckOption {
	return func(handler *checkerHandler) {
		handler.intervalTime = interval
	}
}

func WithTimeout(timeout time.Duration) CheckOption {
	return func(handler *checkerHandler) {
		handler.timeout = timeout
	}
}
