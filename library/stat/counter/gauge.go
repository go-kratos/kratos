package counter

import "sync/atomic"

var _ Counter = new(gaugeCounter)

// A value is a thread-safe counter implementation.
type gaugeCounter int64

// NewGauge return a guage counter.
func NewGauge() Counter {
	return new(gaugeCounter)
}

// Add method increments the counter by some value and return new value
func (v *gaugeCounter) Add(val int64) {
	atomic.AddInt64((*int64)(v), val)
}

// Value method returns the counter's current value.
func (v *gaugeCounter) Value() int64 {
	return atomic.LoadInt64((*int64)(v))
}

// Reset reset the counter.
func (v *gaugeCounter) Reset() {
	atomic.StoreInt64((*int64)(v), 0)
}
