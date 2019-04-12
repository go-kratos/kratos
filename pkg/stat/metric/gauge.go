package metric

import "sync/atomic"

var _ Metric = &gauge{}

// Gauge stores a numerical value that can be add arbitrarily.
type Gauge interface {
	Metric
	// Sets sets the value to the given number.
	Set(int64)
}

// GaugeOpts is an alias of Opts.
type GaugeOpts Opts

type gauge struct {
	val int64
}

// NewGauge creates a new Gauge based on the GaugeOpts.
func NewGauge(opts GaugeOpts) Gauge {
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
