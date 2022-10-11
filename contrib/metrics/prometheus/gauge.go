package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-kratos/kratos/v2/metrics"
)

var _ metrics.Gauge = (*gauge)(nil)

type gauge struct {
	gv  *prometheus.GaugeVec
	lvs []string
}

// NewGauge new a prometheus gauge and returns Gauge.
func NewGauge(gv *prometheus.GaugeVec) metrics.Gauge {
	return &gauge{
		gv: gv,
	}
}

func (g *gauge) With(lvs ...string) metrics.Gauge {
	return &gauge{
		gv:  g.gv,
		lvs: lvs,
	}
}

func (g *gauge) Set(value float64) {
	g.gv.WithLabelValues(g.lvs...).Set(value)
}

func (g *gauge) Add(delta float64) {
	g.gv.WithLabelValues(g.lvs...).Add(delta)
}

func (g *gauge) Sub(delta float64) {
	g.gv.WithLabelValues(g.lvs...).Sub(delta)
}
