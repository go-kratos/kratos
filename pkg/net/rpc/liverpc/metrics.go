package liverpc

import (
	"go-common/library/stat/metric"
)

const namespace = "liverpc_client"

var (
	_metricReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "liverpc client requests duration(ms).",
		Labels:    []string{"method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})
	_metricReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "liverpc client requests code count.",
		Labels:    []string{"method", "code"},
	})
)
