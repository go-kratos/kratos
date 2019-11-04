package hbase

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"

	"github.com/bilibili/kratos/pkg/stat/metric"
)

const namespace = "hbase_client"

var (
	_metricReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "hbase client requests duration(ms).",
		Labels:    []string{"name", "addr", "command"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	})
	_metricReqErr = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "error_total",
		Help:      "mysql client requests error count.",
		Labels:    []string{"name", "addr", "command", "error"},
	})
)

func codeFromErr(err error) string {
	code := "unknown_error"
	switch err {
	case gohbase.ErrClientClosed:
		code = "client_closed"
	case gohbase.ErrCannotFindRegion:
		code = "connot_find_region"
	case gohbase.TableNotFound:
		code = "table_not_found"
		//case gohbase.ErrRegionUnavailable:
		//	code = "region_unavailable"
	}
	return code
}

// MetricsHook if stats is nil use stat.DB as default.
func MetricsHook(config *Config) HookFunc {
	return func(ctx context.Context, call hrpc.Call, customName string) func(err error) {
		now := time.Now()
		if customName == "" {
			customName = call.Name()
		}
		return func(err error) {
			durationMs := int64(time.Since(now) / time.Millisecond)
			_metricReqDur.Observe(durationMs, strings.Join(config.Zookeeper.Addrs, ","), "", customName)
			if err != nil && err != io.EOF {
				_metricReqErr.Inc(strings.Join(config.Zookeeper.Addrs, ","), "", customName, codeFromErr(err))
			}
		}
	}
}
