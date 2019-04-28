package workpool

import (
	"fmt"
	"time"
)

// Task interface
type Task interface {
	Run() *[]byte
}

// FutureTask out must be blocking chan (size=0)
type FutureTask struct {
	T   Task
	out chan *[]byte
}

// NewFutureTask .
func NewFutureTask(t Task) *FutureTask {
	return &FutureTask{
		T:   t,
		out: make(chan *[]byte, 1),
	}
}

// Wait for task return until timeout
func (ft *FutureTask) Wait(timeout time.Duration) (res *[]byte, err error) {
	select {
	case res = <-ft.out:
	case <-time.After(timeout):
		err = fmt.Errorf("task(%+v) timeout", ft)
	}
	return
}
