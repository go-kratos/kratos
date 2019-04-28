// Package timerqueue implements a priority queue for objects scheduled at a
// particular time.
package timerqueue

import (
	"container/heap"
	"errors"
	"time"
)

// Timer is an interface that types implement to schedule and receive OnTimer
// callbacks.
type Timer interface {
	OnTimer(t time.Time)
}

//NewTimerWrapper util struct
func NewTimerWrapper(fun TimerFunc) (result *TimerWrapper) {
	result = &TimerWrapper{fun}
	return
}

//TimerFunc timer function
type TimerFunc func(t time.Time)

//TimerWrapper just a time wrapper
type TimerWrapper struct {
	fun TimerFunc
}

//OnTimer ontimer
func (t *TimerWrapper) OnTimer(tm time.Time) {
	t.fun(tm)
}

// Queue is a time-sorted collection of Timer objects.
type Queue struct {
	heap  timerHeap
	table map[Timer]*timerData
}

type timerData struct {
	timer  Timer
	time   time.Time
	index  int
	period time.Duration // if > 0, this will be a periodically event
}

// New creates a new timer priority queue.
func New() *Queue {
	return &Queue{
		table: make(map[Timer]*timerData),
	}
}

// Len returns the current number of timer objects in the queue.
func (q *Queue) Len() int {
	return len(q.heap)
}

// Schedule schedules a timer for exectuion at time tm. If the
// timer was already scheduled, it is rescheduled.
func (q *Queue) Schedule(t Timer, tm time.Time) {
	q.ScheduleRepeat(t, tm, 0)
}

// ScheduleRepeat give 0 duration, will not be repeatedly event
func (q *Queue) ScheduleRepeat(t Timer, tm time.Time, period time.Duration) {
	if data, ok := q.table[t]; !ok {
		data = &timerData{t, tm, 0, period}
		heap.Push(&q.heap, data)
		q.table[t] = data
	} else {
		data.time = tm
		heap.Fix(&q.heap, data.index)
	}
}

// Unschedule unschedules a timer's execution.
func (q *Queue) Unschedule(t Timer) {
	if data, ok := q.table[t]; ok {
		heap.Remove(&q.heap, data.index)
		delete(q.table, t)
	}
}

// GetTime returns the time at which the timer is scheduled.
// If the timer isn't currently scheduled, an error is returned.
func (q *Queue) GetTime(t Timer) (tm time.Time, err error) {
	if data, ok := q.table[t]; ok {
		return data.time, nil
	}
	return time.Time{}, errors.New("timerqueue: timer not scheduled")
}

// IsScheduled returns true if the timer is currently scheduled.
func (q *Queue) IsScheduled(t Timer) bool {
	_, ok := q.table[t]
	return ok
}

// Clear unschedules all currently scheduled timers.
func (q *Queue) Clear() {
	q.heap, q.table = nil, make(map[Timer]*timerData)
}

// PopFirst removes and returns the next timer to be scheduled and
// the time at which it is scheduled to run.
func (q *Queue) PopFirst() (t Timer, tm time.Time) {
	if len(q.heap) > 0 {
		data := heap.Pop(&q.heap).(*timerData)
		delete(q.table, data.timer)
		return data.timer, data.time
	}
	return nil, time.Time{}
}

// PeekFirst returns the next timer to be scheduled and the time
// at which it is scheduled to run. It does not modify the contents
// of the timer queue.
func (q *Queue) PeekFirst() (t Timer, tm time.Time) {
	if len(q.heap) > 0 {
		return q.heap[0].timer, q.heap[0].time
	}
	return nil, time.Time{}
}

// Advance executes OnTimer callbacks for all timers scheduled to be
// run before the time 'tm'. Executed timers are removed from the
// timer queue.
func (q *Queue) Advance(tm time.Time) {
	for len(q.heap) > 0 && !tm.Before(q.heap[0].time) {
		data := q.heap[0]
		heap.Remove(&q.heap, data.index)
		if data.period > 0 {
			data.time = data.time.Add(data.period)
			heap.Push(&q.heap, data)
		} else {
			delete(q.table, data.timer)
		}
		data.timer.OnTimer(data.time)
	}
}

/*
 * timerHeap
 */

type timerHeap []*timerData

//Len len interface
func (h timerHeap) Len() int {
	return len(h)
}

//Less less interface
func (h timerHeap) Less(i, j int) bool {
	return h[i].time.Before(h[j].time)
}

//Swap swap interface
func (h timerHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index, h[j].index = i, j
}

//Push push interface
func (h *timerHeap) Push(x interface{}) {
	data := x.(*timerData)
	*h = append(*h, data)
	data.index = len(*h) - 1
}

//Pop pop interface
func (h *timerHeap) Pop() interface{} {
	n := len(*h)
	data := (*h)[n-1]
	*h = (*h)[:n-1]
	data.index = -1
	return data
}
