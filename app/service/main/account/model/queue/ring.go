/*
Copyright 2014 Workiva, LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package queue

import (
	"runtime"
	"sync/atomic"
	"time"
)

// roundUp takes a uint64 greater than 0 and rounds it up to the next
// power of 2.
func roundUp(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}

type node struct {
	position uint64
	data     interface{}
}

type nodes []*node

// RingBuffer is a MPMC buffer that achieves threadsafety with CAS operations
// only.  A put on full or get on empty call will block until an item
// is put or retrieved.  Calling Dispose on the RingBuffer will unblock
// any blocked threads with an error.  This buffer is similar to the buffer
// described here: http://www.1024cores.net/home/lock-free-algorithms/queues/bounded-mpmc-queue
// with some minor additions.
type RingBuffer struct {
	_padding0      [8]uint64
	queue          uint64
	_padding1      [8]uint64
	dequeue        uint64
	_padding2      [8]uint64
	mask, disposed uint64
	_padding3      [8]uint64
	nodes          nodes
}

func (rb *RingBuffer) init(size uint64) {
	size = roundUp(size)
	rb.nodes = make(nodes, size)
	for i := uint64(0); i < size; i++ {
		rb.nodes[i] = &node{position: i}
	}
	rb.mask = size - 1 // so we don't have to do this with every put/get operation
}

// Put adds the provided item to the queue.  If the queue is full, this
// call will block until an item is added to the queue or Dispose is called
// on the queue.  An error will be returned if the queue is disposed.
func (rb *RingBuffer) Put(item interface{}) error {
	_, err := rb.put(item, false)
	return err
}

// Offer adds the provided item to the queue if there is space.  If the queue
// is full, this call will return false.  An error will be returned if the
// queue is disposed.
func (rb *RingBuffer) Offer(item interface{}) (bool, error) {
	return rb.put(item, true)
}

func (rb *RingBuffer) put(item interface{}, offer bool) (bool, error) {
	var n *node
	pos := atomic.LoadUint64(&rb.queue)
L:
	for {
		if atomic.LoadUint64(&rb.disposed) == 1 {
			return false, ErrDisposed
		}

		n = rb.nodes[pos&rb.mask]
		seq := atomic.LoadUint64(&n.position)
		switch dif := seq - pos; {
		case dif == 0:
			if atomic.CompareAndSwapUint64(&rb.queue, pos, pos+1) {
				break L
			}
		// case dif < 0:
		// 	panic(`Ring buffer in a compromised state during a put operation.`)
		default:
			pos = atomic.LoadUint64(&rb.queue)
		}

		if offer {
			return false, nil
		}

		runtime.Gosched() // free up the cpu before the next iteration
	}

	n.data = item
	atomic.StoreUint64(&n.position, pos+1)
	return true, nil
}

// Get will return the next item in the queue.  This call will block
// if the queue is empty.  This call will unblock when an item is added
// to the queue or Dispose is called on the queue.  An error will be returned
// if the queue is disposed.
func (rb *RingBuffer) Get() (interface{}, error) {
	return rb.Poll(0)
}

// Poll will return the next item in the queue.  This call will block
// if the queue is empty.  This call will unblock when an item is added
// to the queue, Dispose is called on the queue, or the timeout is reached. An
// error will be returned if the queue is disposed or a timeout occurs. A
// non-positive timeout will block indefinitely.
func (rb *RingBuffer) Poll(timeout time.Duration) (interface{}, error) {
	var (
		n     *node
		pos   = atomic.LoadUint64(&rb.dequeue)
		start time.Time
	)
	if timeout > 0 {
		start = time.Now()
	}
L:
	for {
		if atomic.LoadUint64(&rb.disposed) == 1 {
			return nil, ErrDisposed
		}

		n = rb.nodes[pos&rb.mask]
		seq := atomic.LoadUint64(&n.position)
		switch dif := seq - (pos + 1); {
		case dif == 0:
			if atomic.CompareAndSwapUint64(&rb.dequeue, pos, pos+1) {
				break L
			}
		// case dif < 0:
		// 	panic(`Ring buffer in compromised state during a get operation.`)
		default:
			pos = atomic.LoadUint64(&rb.dequeue)
		}

		if timeout > 0 && time.Since(start) >= timeout {
			return nil, ErrTimeout
		}

		runtime.Gosched() // free up the cpu before the next iteration
	}
	data := n.data
	n.data = nil
	atomic.StoreUint64(&n.position, pos+rb.mask+1)
	return data, nil
}

// Len returns the number of items in the queue.
func (rb *RingBuffer) Len() uint64 {
	return atomic.LoadUint64(&rb.queue) - atomic.LoadUint64(&rb.dequeue)
}

// Cap returns the capacity of this ring buffer.
func (rb *RingBuffer) Cap() uint64 {
	return uint64(len(rb.nodes))
}

// Dispose will dispose of this queue and free any blocked threads
// in the Put and/or Get methods.  Calling those methods on a disposed
// queue will return an error.
func (rb *RingBuffer) Dispose() {
	atomic.CompareAndSwapUint64(&rb.disposed, 0, 1)
}

// IsDisposed will return a bool indicating if this queue has been
// disposed.
func (rb *RingBuffer) IsDisposed() bool {
	return atomic.LoadUint64(&rb.disposed) == 1
}

// NewRingBuffer will allocate, initialize, and return a ring buffer
// with the specified size.
func NewRingBuffer(size uint64) *RingBuffer {
	rb := &RingBuffer{}
	rb.init(size)
	return rb
}
