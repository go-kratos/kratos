package metric

import "time"

var _ Metric = &rollingGauge{}
var _ Aggregation = &rollingGauge{}

// RollingGauge represents a ring window based on time duration.
// e.g. [[1, 2], [1, 2, 3], [1,2, 3, 4]]
type RollingGauge interface {
	Metric
	Aggregation
	// Reduce applies the reduction function to all buckets within the window.
	Reduce(func(Iterator) float64) float64
}

// RollingGaugeOpts contains the arguments for creating RollingGauge.
type RollingGaugeOpts struct {
	Size           int
	BucketDuration time.Duration
}

type rollingGauge struct {
	policy *RollingPolicy
}

// NewRollingGauge creates a new RollingGauge baseed on RollingGaugeOpts.
func NewRollingGauge(opts RollingGaugeOpts) RollingGauge {
	window := NewWindow(WindowOpts{Size: opts.Size})
	policy := NewRollingPolicy(window, RollingPolicyOpts{BucketDuration: opts.BucketDuration})
	return &rollingGauge{
		policy: policy,
	}
}

func (r *rollingGauge) Add(val int64) {
	r.policy.Append(float64(val))
}

func (r *rollingGauge) Reduce(f func(Iterator) float64) float64 {
	return r.policy.Reduce(f)
}

func (r *rollingGauge) Avg() float64 {
	return r.policy.Reduce(Avg)
}

func (r *rollingGauge) Min() float64 {
	return r.policy.Reduce(Min)
}

func (r *rollingGauge) Max() float64 {
	return r.policy.Reduce(Max)
}

func (r *rollingGauge) Sum() float64 {
	return r.policy.Reduce(Sum)
}

func (r *rollingGauge) Value() int64 {
	return int64(r.Sum())
}
