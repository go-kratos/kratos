package hbase

import (
	"context"
	"io"
	"time"

	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"

	"github.com/bilibili/kratos/pkg/stat"
)

func codeFromErr(err error) string {
	code := "unknown_error"
	switch err {
	case gohbase.ErrClientClosed:
		code = "client_closed"
	case gohbase.ErrConnotFindRegion:
		code = "connot_find_region"
	case gohbase.TableNotFound:
		code = "table_not_found"
	case gohbase.ErrRegionUnavailable:
		code = "region_unavailable"
	}
	return code
}

// MetricsHook if stats is nil use stat.DB as default.
func MetricsHook(stats stat.Stat) HookFunc {
	if stats == nil {
		stats = stat.DB
	}
	return func(ctx context.Context, call hrpc.Call, customName string) func(err error) {
		now := time.Now()
		if customName == "" {
			customName = call.Name()
		}
		method := "hbase:" + customName
		return func(err error) {
			durationMs := int64(time.Since(now) / time.Millisecond)
			stats.Timing(method, durationMs)
			if err != nil && err != io.EOF {
				stats.Incr(method, codeFromErr(err))
			}
		}
	}
}
