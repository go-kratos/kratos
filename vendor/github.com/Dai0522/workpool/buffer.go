package workpool

import (
	"errors"
	"runtime"
	"sync/atomic"
)

// ringBuffer .
type ringBuffer struct {
	capacity       uint64
	mask           uint64
	padding1       [7]uint64
	lastCommintIdx uint64
	padding2       [7]uint64
	nextFreeIdx    uint64
	padding3       [7]uint64
	readerIdx      uint64
	padding4       [7]uint64
	slots          []*worker
}

// newRingBuffer .
func newRingBuffer(c uint64) (*ringBuffer, error) {
	if c == 0 || c&3 != 0 {
		return nil, errors.New("capacity must be N power of 2")
	}
	return &ringBuffer{
		lastCommintIdx: 0,
		nextFreeIdx:    1,
		readerIdx:      0,
		capacity:       c,
		mask:           c - 1,
		slots:          make([]*worker, c),
	}, nil
}

// push .
func (r *ringBuffer) push(w *worker) error {
	var head, tail, next uint64
	for {
		head = r.nextFreeIdx
		tail = r.readerIdx
		if (head > tail+r.capacity-2) || (head < tail-1) {
			return errors.New("buffer is full")
		}

		next = (head + 1) & r.mask
		if atomic.CompareAndSwapUint64(&r.nextFreeIdx, head, next) {
			break
		}
		runtime.Gosched()
	}
	r.slots[head] = w

	for !atomic.CompareAndSwapUint64(&r.lastCommintIdx, head-1, head) {
		runtime.Gosched()
	}
	return nil
}

// pop .
func (r *ringBuffer) pop() *worker {
	var head, next uint64
	for {
		head = r.readerIdx
		if head == r.lastCommintIdx {
			return r.slots[head]
		}
		next = (head + 1) & r.mask
		if atomic.CompareAndSwapUint64(&r.readerIdx, head, next) {
			break
		}
		runtime.Gosched()
	}
	return r.slots[head]
}

// size .
func (r *ringBuffer) size() uint64 {
	return r.lastCommintIdx - r.readerIdx
}
