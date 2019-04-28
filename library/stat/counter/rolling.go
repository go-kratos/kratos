package counter

import (
	"sync"
	"time"
)

type bucket struct {
	val  int64
	next *bucket
}

func (b *bucket) Add(val int64) {
	b.val += val
}

func (b *bucket) Value() int64 {
	return b.val
}

func (b *bucket) Reset() {
	b.val = 0
}

var _ Counter = new(rollingCounter)

type rollingCounter struct {
	mu         sync.RWMutex
	buckets    []bucket
	bucketTime int64
	lastAccess int64
	cur        *bucket
}

// NewRolling creates a new window. windowTime is the time covering the entire
// window. windowBuckets is the number of buckets the window is divided into.
// An example: a 10 second window with 10 buckets will have 10 buckets covering
// 1 second each.
func NewRolling(window time.Duration, winBucket int) Counter {
	buckets := make([]bucket, winBucket)
	bucket := &buckets[0]
	for i := 1; i < winBucket; i++ {
		bucket.next = &buckets[i]
		bucket = bucket.next
	}
	bucket.next = &buckets[0]
	bucketTime := time.Duration(window.Nanoseconds() / int64(winBucket))
	return &rollingCounter{
		cur:        &buckets[0],
		buckets:    buckets,
		bucketTime: int64(bucketTime),
		lastAccess: time.Now().UnixNano(),
	}
}

// Add increments the counter by value and return new value.
func (r *rollingCounter) Add(val int64) {
	r.mu.Lock()
	r.lastBucket().Add(val)
	r.mu.Unlock()
}

// Value get the counter value.
func (r *rollingCounter) Value() (sum int64) {
	now := time.Now().UnixNano()
	r.mu.RLock()
	b := r.cur
	i := r.elapsed(now)
	for j := 0; j < len(r.buckets); j++ {
		// skip all future reset bucket.
		if i > 0 {
			i--
		} else {
			sum += b.Value()
		}
		b = b.next
	}
	r.mu.RUnlock()
	return
}

//  Reset reset the counter.
func (r *rollingCounter) Reset() {
	r.mu.Lock()
	for i := range r.buckets {
		r.buckets[i].Reset()
	}
	r.mu.Unlock()
}

func (r *rollingCounter) elapsed(now int64) (i int) {
	var e int64
	if e = now - r.lastAccess; e <= r.bucketTime {
		return
	}
	if i = int(e / r.bucketTime); i > len(r.buckets) {
		i = len(r.buckets)
	}
	return
}

func (r *rollingCounter) lastBucket() (b *bucket) {
	now := time.Now().UnixNano()
	b = r.cur
	// reset the buckets between now and number of buckets ago. If
	// that is more that the existing buckets, reset all.
	if i := r.elapsed(now); i > 0 {
		r.lastAccess = now
		for ; i > 0; i-- {
			// replace the next used bucket.
			b = b.next
			b.Reset()
		}
	}
	r.cur = b
	return
}
