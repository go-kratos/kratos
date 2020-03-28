package sql

import "github.com/go-kratos/kratos/pkg/stat/metric"

const namespace = "mysql_client"

var (
	_metricReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "mysql client requests duration(ms).",
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
	_metricConnTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "connections",
		Name:      "total",
		Help:      "mysql client connections total count.",
		Labels:    []string{"name", "addr", "state"},
	})
	_metricConnCurrent = metric.NewGaugeVec(&metric.GaugeVecOpts{
		Namespace: namespace,
		Subsystem: "connections",
		Name:      "current",
		Help:      "mysql client connections current.",
		Labels:    []string{"name", "addr", "state"},
	})
)
