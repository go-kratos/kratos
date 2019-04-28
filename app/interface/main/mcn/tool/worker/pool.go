package worker

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"go-common/library/log"
)

const (
	_ratio = float32(0.8)
)

var (
	_default = &Conf{
		QueueSize:     1024,
		WorkerProcMax: 32,
		WorkerNumber:  runtime.NumCPU() - 1,
	}
)

// Conf .
type Conf struct {
	QueueSize     int
	WorkerProcMax int
	WorkerNumber  int
}

// Pool .
type Pool struct {
	c            *Conf
	queue        chan func()
	workerNumber int
	close        chan struct{}
	wg           sync.WaitGroup
}

// New .
func New(conf *Conf) (w *Pool) {
	if conf == nil {
		conf = _default
	}
	w = &Pool{
		c:            conf,
		queue:        make(chan func(), conf.QueueSize),
		workerNumber: conf.WorkerNumber,
		close:        make(chan struct{}),
	}
	w.start()
	go w.moni()
	return
}

func (w *Pool) start() {
	for i := 0; i < w.workerNumber; i++ {
		w.wg.Add(1)
		go w.workerRoutine()
	}
}

func (w *Pool) moni() {
	var conf = w.c
	for {
		time.Sleep(time.Second * 5)
		var ratio = float32(len(w.queue)) / float32(conf.QueueSize)
		if ratio >= _ratio {
			if w.workerNumber >= conf.WorkerProcMax {
				log.Warn("work thread more than max(%d)", conf.WorkerProcMax)
				return
			}
			var next = minInt(w.workerNumber<<1, w.c.WorkerProcMax)
			var diff = next - w.workerNumber
			log.Info("current thread count=%d, queue ratio=%f, create new thread number=(%d)", w.workerNumber, ratio, diff)
			for i := 0; i < diff; i++ {
				w.wg.Add(1)
				go w.workerRoutine()
			}
			w.workerNumber = next
		}
	}
}

// Close .
func (w *Pool) Close() {
	close(w.close)
}

// Wait .
func (w *Pool) Wait() {
	w.wg.Wait()
}

func (w *Pool) workerRoutine() {
	defer func() {
		w.wg.Done()
		if x := recover(); x != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Error("w.workerRoutine panic(%+v) :\n %s", x, buf)
			w.wg.Add(1)
			go w.workerRoutine()
		}
	}()
loop:
	for {
		select {
		case f := <-w.queue:
			f()
		case <-w.close:
			log.Info("workerRoutine close()")
			break loop
		}
	}
	for f := range w.queue {
		f()
	}
}

// Add .
func (w *Pool) Add(f func()) error {
	select {
	case w.queue <- f:
	default:
		return fmt.Errorf("task channel is full")
	}
	return nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
