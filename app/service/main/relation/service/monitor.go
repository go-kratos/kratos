package service

import (
	"context"
	"time"
)

// Monitor if mid is monitored
func (s *Service) Monitor(c context.Context, mid int64) (monitor bool, err error) {
	if !s.c.Relation.Monitor {
		return
	}
	return s.dao.MonitorCache(c, mid)
}

// AddMonitor add mid to monitor table.
func (s *Service) AddMonitor(c context.Context, mid int64) (err error) {
	if _, err = s.dao.AddMonitor(c, mid, time.Now()); err != nil {
		return
	}
	return s.dao.SetMonitorCache(c, mid)
}

// DelMonitor del mid from monitor table
func (s *Service) DelMonitor(c context.Context, mid int64) (err error) {
	if _, err = s.dao.DelMonitor(c, mid); err != nil {
		return
	}
	return s.dao.DelMonitorCache(c, mid)
}

// LoadMonitor load monitor
func (s *Service) LoadMonitor(c context.Context) (err error) {
	var (
		mids []int64
	)
	if mids, _ = s.dao.LoadMonitor(c); mids != nil {
		return s.dao.LoadMonitorCache(c, mids)
	}
	return
}
