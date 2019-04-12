package metric

var _ Metric = &pointGauge{}
var _ Aggregation = &pointGauge{}

// PointGauge represents a ring window.
// Every buckets within the window contains one point.
// When the window is full, the earliest point will be overwrite.
type PointGauge interface {
	Aggregation
	Metric
	// Reduce applies the reduction function to all buckets within the window.
	Reduce(func(Iterator) float64) float64
}

// PointGaugeOpts contains the arguments for creating PointGauge.
type PointGaugeOpts struct {
	// Size represents the bucket size within the window.
	Size int
}

type pointGauge struct {
	policy *PointPolicy
}

// NewPointGauge creates a new PointGauge based on PointGaugeOpts.
func NewPointGauge(opts PointGaugeOpts) PointGauge {
	window := NewWindow(WindowOpts{Size: opts.Size})
	policy := NewPointPolicy(window)
	return &pointGauge{
		policy: policy,
	}
}

func (r *pointGauge) Add(val int64) {
	r.policy.Append(float64(val))
}

func (r *pointGauge) Reduce(f func(Iterator) float64) float64 {
	return r.policy.Reduce(f)
}

func (r *pointGauge) Avg() float64 {
	return r.policy.Reduce(Avg)
}

func (r *pointGauge) Min() float64 {
	return r.policy.Reduce(Min)
}

func (r *pointGauge) Max() float64 {
	return r.policy.Reduce(Max)
}

func (r *pointGauge) Sum() float64 {
	return r.policy.Reduce(Sum)
}

func (r *pointGauge) Value() int64 {
	return int64(r.Sum())
}
