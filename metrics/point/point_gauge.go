package point

import (
	"github.com/go-kratos/kratos/v2/metrics"
	"github.com/go-kratos/kratos/v2/metrics/rolling"
)

var _ metrics.Metric = &pointGauge{}
var _ metrics.Aggregation = &pointGauge{}

// PointGauge represents a ring window.
// Every buckets within the window contains one point.
// When the window is full, the earliest point will be overwrite.
type PointGauge interface {
	metrics.Aggregation
	metrics.Metric
	// Reduce applies the reduction function to all buckets within the window.
	Reduce(func(rolling.Iterator) float64) float64
}

// PointGaugeOpts contains the arguments for creating PointGauge.
type PointGaugeOpts struct {
	// Size represents the bucket size within the window.
	Size int
}

type pointGauge struct {
	policy *rolling.PointPolicy
}

// NewPointGauge creates a new PointGauge based on PointGaugeOpts.
func NewPointGauge(opts PointGaugeOpts) PointGauge {
	window := rolling.NewWindow(rolling.WindowOpts{Size: opts.Size})
	policy := rolling.NewPointPolicy(window)
	return &pointGauge{
		policy: policy,
	}
}

func (r *pointGauge) Add(val int64) {
	r.policy.Append(float64(val))
}

func (r *pointGauge) Reduce(f func(rolling.Iterator) float64) float64 {
	return r.policy.Reduce(f)
}

func (r *pointGauge) Avg() float64 {
	return r.policy.Reduce(rolling.Avg)
}

func (r *pointGauge) Min() float64 {
	return r.policy.Reduce(rolling.Min)
}

func (r *pointGauge) Max() float64 {
	return r.policy.Reduce(rolling.Max)
}

func (r *pointGauge) Sum() float64 {
	return r.policy.Reduce(rolling.Sum)
}

func (r *pointGauge) Value() int64 {
	return int64(r.Sum())
}
