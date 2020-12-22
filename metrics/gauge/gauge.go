package gauge

import (
	"sync/atomic"

	"github.com/go-kratos/kratos/v2/metrics"
)

var _ metrics.Metric = &gauge{}

// GaugeOpts is an alias of Opts.
type GaugeOpts metrics.Opts

type gauge struct {
	val int64
}

// NewGauge creates a new Gauge based on the GaugeOpts.
func NewGauge(opts GaugeOpts) metrics.Gauge {
	return &gauge{}
}

func (g *gauge) Add(val int64) {
	atomic.AddInt64(&g.val, val)
}

func (g *gauge) Set(val int64) {
	old := atomic.LoadInt64(&g.val)
	atomic.CompareAndSwapInt64(&g.val, old, val)
}

func (g *gauge) Value() int64 {
	return atomic.LoadInt64(&g.val)
}

// GaugeVecOpts is an alias of VectorOpts.
type GaugeVecOpts metrics.VectorOpts
