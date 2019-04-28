package service

import (
	"context"
	"go-common/app/interface/openplatform/article/conf"
)

var _sentinel = &conf.Sentinel{
	EnableSentinel:     1,
	DurationSample:     100,
	MonitorCountSample: 100,
	MonitorRateSample:  100,
	DebugSample:        100,
}

// Sentinel .
func (s *Service) Sentinel(c context.Context) *conf.Sentinel {
	if s.c.Sentinel == nil {
		return _sentinel
	}
	return s.c.Sentinel
}
