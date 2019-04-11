package hbase

import (
	"context"
	"time"

	"github.com/tsuna/gohbase/hrpc"

	"github.com/bilibili/kratos/pkg/log"
)

// NewSlowLogHook log slow operation.
func NewSlowLogHook(threshold time.Duration) HookFunc {
	return func(ctx context.Context, call hrpc.Call, customName string) func(err error) {
		start := time.Now()
		return func(error) {
			duration := time.Since(start)
			if duration < threshold {
				return
			}
			log.Warn("hbase slow log: %s %s %s time: %s", customName, call.Table(), call.Key(), duration)
		}
	}
}
