package bbr

import (
	"sync"
	"time"
)

// RollingPolicy is a policy for ring window based on time duration.
// RollingPolicy moves bucket offset with time duration.
// e.g. If the last point is appended one bucket duration ago,
// RollingPolicy will increment current offset.
type RollingPolicy struct {
	mu     sync.RWMutex
	size   int
	window *Window
	offset int

	bucketDuration time.Duration
	lastAppendTime time.Time
}

// RollingPolicyOpts contains the arguments for creating RollingPolicy.
type RollingPolicyOpts struct {
	BucketDuration time.Duration
}

// NewRollingPolicy creates a new RollingPolicy based on the given window and RollingPolicyOpts.
func NewRollingPolicy(window *Window, opts RollingPolicyOpts) *RollingPolicy {
	return &RollingPolicy{
		window: window,
		size:   window.Size(),
		offset: 0,

		bucketDuration: opts.BucketDuration,
		lastAppendTime: time.Now(),
	}
}

// timespan returns passed bucket number since lastAppendTime,
// if it is one bucket duration earlier than the last recorded
// time, it will return the size.
func (r *RollingPolicy) timespan() int {
	v := int(time.Since(r.lastAppendTime) / r.bucketDuration)
	if v > -1 { // maybe time backwards
		return v
	}
	return r.size
}

// apply applies function f with value val on
// current offset bucket, expired bucket will be reset
func (r *RollingPolicy) apply(f func(offset int, val float64), val float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// calculate current offset
	timespan := r.timespan()
	if timespan > 0 {
		start := (r.offset + 1) % r.size
		end := (r.offset + timespan) % r.size
		if timespan > r.size {
			timespan = r.size
		}
		// reset the expired buckets
		r.window.ResetBuckets(start, timespan)
		r.offset = end
		r.lastAppendTime = r.lastAppendTime.Add(time.Duration(timespan * int(r.bucketDuration)))
	}
	f(r.offset, val)
}

// Append appends the given points to the window.
func (r *RollingPolicy) Append(val float64) {
	r.apply(r.window.Append, val)
}

// Add adds the given value to the latest point within bucket.
func (r *RollingPolicy) Add(val float64) {
	r.apply(r.window.Add, val)
}

// Reduce applies the reduction function to all buckets within the window.
func (r *RollingPolicy) Reduce(f func(Iterator) float64) (val float64) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	timespan := r.timespan()
	if count := r.size - timespan; count > 0 {
		offset := r.offset + timespan + 1
		if offset >= r.size {
			offset = offset - r.size
		}
		val = f(r.window.Iterator(offset, count))
	}
	return val
}
