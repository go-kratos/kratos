package health

import (
	"golang.org/x/net/context"
)

type CheckerMgr struct {
	checkers map[string]checker
	ctx      context.Context
	cancel   func()
	watchers []chan string
}

func New(ctx context.Context) *CheckerMgr {
	c, cancel := context.WithCancel(ctx)
	return &CheckerMgr{
		checkers: make(map[string]checker),
		ctx:      c,
		cancel:   cancel,
	}
}

func (c *CheckerMgr) Start() {
	for _, v := range c.checkers {
		go func() {
			v.run(c.ctx)
		}()
	}
}

func (c *CheckerMgr) Stop() {
	c.cancel()
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
