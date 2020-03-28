package blademaster

import "github.com/go-kratos/kratos/pkg/stat/metric"

const (
	clientNamespace = "http_client"
	serverNamespace = "http_server"
)

var (
	_metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Labels:    []string{"path", "caller", "method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})
	_metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
		Labels:    []string{"path", "caller", "method", "code"},
	})
	_metricServerBBR = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "",
		Name:      "bbr_total",
		Help:      "http server bbr total.",
		Labels:    []string{"url", "method"},
	})
	_metricClientReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http client requests duration(ms).",
		Labels:    []string{"path", "method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})
	_metricClientReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http client requests code count.",
		Labels:    []string{"path", "method", "code"},
	})
)
