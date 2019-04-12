package metric

import "sync"

// PointPolicy is a policy of points within the window.
// PointPolicy wraps the window and make it seem like ring-buf.
// When using PointPolicy, every buckets within the windows contains at more one point.
// e.g. [[1], [2], [3]]
type PointPolicy struct {
	mu     sync.RWMutex
	size   int
	window *Window
	offset int
}

// NewPointPolicy creates a new PointPolicy.
func NewPointPolicy(window *Window) *PointPolicy {
	return &PointPolicy{
		window: window,
		size:   window.Size(),
		offset: -1,
	}
}

func (p *PointPolicy) prevOffset() int {
	return p.offset
}

func (p *PointPolicy) nextOffset() int {
	return (p.prevOffset() + 1) % p.size
}

func (p *PointPolicy) updateOffset(offset int) {
	p.offset = offset
}

// Append appends the given points to the window.
func (p *PointPolicy) Append(val float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	offset := p.nextOffset()
	p.window.ResetBucket(offset)
	p.window.Append(offset, val)
	p.updateOffset(offset)
}

// Reduce applies the reduction function to all buckets within the window.
func (p *PointPolicy) Reduce(f func(Iterator) float64) float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	offset := p.offset + 1
	if offset == p.size {
		offset = 0
	}
	iterator := p.window.Iterator(offset, p.size)
	return f(iterator)
}
