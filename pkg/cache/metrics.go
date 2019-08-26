package cache

import "github.com/bilibili/kratos/pkg/stat/metric"

const _metricNamespace = "cache"

// be used in tool/kratos-gen-bts
var (
	MetricHits = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: _metricNamespace,
		Subsystem: "",
		Name:      "hits_total",
		Help:      "cache hits total.",
		Labels:    []string{"name"},
	})
	MetricMisses = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: _metricNamespace,
		Subsystem: "",
		Name:      "misses_total",
		Help:      "cache misses total.",
		Labels:    []string{"name"},
	})
)
