package trace

import (
	"github.com/bilibili/kratos/pkg/stat/metric"
)

const namespace = "dapper_client"

var (
	jaegerReportProtocolCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "report",
		Name:      "protocol",
		Help:      "dapper client report protocol count",
		Labels:    []string{"report"},
	})
	jaegerReportDroppedCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "report",
		Name:      "dropped",
		Help:      "dapper client dropped span count",
		Labels:    []string{"report"},
	})
	jaegerReportErrorCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: namespace,
		Subsystem: "report",
		Name:      "error",
		Help:      "dapper client send span error count",
		Labels:    []string{"report"},
	})
)
