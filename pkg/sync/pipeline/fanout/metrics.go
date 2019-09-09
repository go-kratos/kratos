package fanout

import (
	"github.com/bilibili/kratos/pkg/stat/metric"
)

const (
	_metricNamespace = "sync"
	_metricSubSystem = "pipeline_fanout"
)

var (
	_metricChanSize = metric.NewGaugeVec(&metric.GaugeVecOpts{
		Namespace: _metricNamespace,
		Subsystem: _metricSubSystem,
		Name:      "chan_len",
		Help:      "sync pipeline fanout current channel size.",
		Labels:    []string{"name"},
	})
	_metricCount = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: _metricNamespace,
		Subsystem: _metricSubSystem,
		Name:      "process_count",
		Help:      "process count",
		Labels:    []string{"name"},
	})
)
