package context

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type mergeCtx struct {
	parent1, parent2 context.Context

	done     chan struct{}
	doneMark uint32
	doneOnce sync.Once
	doneErr  error

	cancelCh   chan struct{}
	cancelOnce sync.Once
}

// Merge merges two contexts into one.
func Merge(parent1, parent2 context.Context) (context.Context, context.CancelFunc) {
	mc := &mergeCtx{
		parent1:  parent1,
		parent2:  parent2,
		done:     make(chan struct{}),
		cancelCh: make(chan struct{}),
	}
	select {
	case <-parent1.Done():
		mc.finish(parent1.Err())
	case <-parent2.Done():
		mc.finish(parent2.Err())
	default:
		go mc.wait()
	}
	return mc, mc.cancel
}

func (mc *mergeCtx) finish(err error) error {
	mc.doneOnce.Do(func() {
		mc.doneErr = err
		atomic.StoreUint32(&mc.doneMark, 1)
		close(mc.done)
	})
	return mc.doneErr
}

func (mc *mergeCtx) wait() {
	var err error
	select {
	case <-mc.parent1.Done():
		err = mc.parent1.Err()
	case <-mc.parent2.Done():
		err = mc.parent2.Err()
	case <-mc.cancelCh:
		err = context.Canceled
	}
	mc.finish(err)
}

func (mc *mergeCtx) cancel() {
	mc.cancelOnce.Do(func() {
		close(mc.cancelCh)
	})
}

// Done implements context.Context.
func (mc *mergeCtx) Done() <-chan struct{} {
	return mc.done
}

// Err implements context.Context.
func (mc *mergeCtx) Err() error {
	if atomic.LoadUint32(&mc.doneMark) != 0 {
		return mc.doneErr
	}
	var err error
	select {
	case <-mc.parent1.Done():
		err = mc.parent1.Err()
	case <-mc.parent2.Done():
		err = mc.parent2.Err()
	case <-mc.cancelCh:
		err = context.Canceled
	default:
		return nil
	}
	return mc.finish(err)
}

// Deadline implements context.Context.
func (mc *mergeCtx) Deadline() (time.Time, bool) {
	d1, ok1 := mc.parent1.Deadline()
	d2, ok2 := mc.parent2.Deadline()
	switch {
	case !ok1:
		return d2, ok2
	case !ok2:
		return d1, ok1
	case d1.Before(d2):
		return d1, true
	default:
		return d2, true
	}
}

// Value implements context.Context.
func (mc *mergeCtx) Value(key interface{}) interface{} {
	if v := mc.parent1.Value(key); v != nil {
		return v
	}
	return mc.parent2.Value(key)
}
