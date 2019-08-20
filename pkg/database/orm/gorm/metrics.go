package gorm

import (
	"time"

	"github.com/bilibili/kratos/pkg/stat/metric"
)

const namespace = "mysql_client_gorm"

var _incConnMetric = func() func(string, string, string) {
	_me := metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "connections",
		Name:      "total",
		Help:      "mysql client connections total count.",
		Labels:    []string{"name", "addr", "state"},
	})
	return func(name string, addr string, state string) {
		_me.Inc(name, addr, state)
	}
}()

var _observeConn = func() func(time.Duration, string, string) {
	_me := metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: namespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "mysql client requests duration(ms).",
		Labels:    []string{"name", "addr"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500},
	})
	return func(dur time.Duration, name string, addr string) {
		_me.Observe(int64(dur/time.Millisecond), name, addr)
	}
}()
