package workpool

import (
	"errors"
	"runtime"
	"sync"
	"time"
)

const (
	stateCreate   = 0
	stateRunning  = 1
	stateStopping = 2
	stateShutdown = 3
)

// PoolConfig .
type PoolConfig struct {
	MaxWorkers     uint64
	MaxIdleWorkers uint64
	MinIdleWorkers uint64
	KeepAlive      time.Duration
}

// Pool .
type Pool struct {
	conf       *PoolConfig
	padding1   [8]uint64
	ready      *ringBuffer
	curWorkers uint64
	padding2   [8]uint64
	lock       sync.Mutex
	state      uint8
	stop       chan uint8
}

// worker .
type worker struct {
	id          uint64
	lastUseTime time.Time
	ftch        chan *FutureTask
}

var wChanCap = func() int {
	// Use blocking worker if GOMAXPROCS=1.
	// This immediately switches Serve to WorkerFunc, which results
	// in higher performance (under go1.5 at least).
	if runtime.GOMAXPROCS(0) == 1 {
		return 0
	}

	// Use non-blocking worker if GOMAXPROCS>1,
	// since otherwise the Serve caller (Acceptor) may lag accepting
	// new task if WorkerFunc is CPU-bound.
	return 1
}()

func newWorker(wid uint64) *worker {
	return &worker{
		id:          wid,
		lastUseTime: time.Now(),
		ftch:        make(chan *FutureTask, wChanCap),
	}
}

// NewWorkerPool .
func NewWorkerPool(capacity uint64, conf *PoolConfig) (p *Pool, err error) {
	if capacity == 0 || capacity&3 != 0 {
		err = errors.New("capacity must bigger than zero and N power of 2")
		return
	}

	rb, err := newRingBuffer(capacity)
	if err != nil {
		return
	}
	p = &Pool{
		conf:       conf,
		ready:      rb,
		curWorkers: 0,
		state:      stateCreate,
		stop:       make(chan uint8, 1),
	}
	return
}

func (p *Pool) changeState(old, new uint8) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.state != old {
		return false
	}

	p.state = new
	return true
}

// Start .
func (p *Pool) Start() error {
	if !p.changeState(stateCreate, stateRunning) {
		return errors.New("workerpool already started")
	}
	go func() {
		defer close(p.stop)
		for {
			p.clean()
			select {
			case <-p.stop:
				p.cleanAll()
				for !p.changeState(stateStopping, stateShutdown) {
					runtime.Gosched()
				}
				return
			default:
				time.Sleep(p.conf.KeepAlive)
			}
		}
	}()
	return nil
}

// Stop .
func (p *Pool) Stop() error {
	if !p.changeState(stateRunning, stateStopping) {
		return errors.New("workerpool is stopping")
	}
	p.stop <- stateStopping
	return nil
}

// Submit .
func (p *Pool) Submit(ft *FutureTask) error {
	w, err := p.getReadyWorker()
	if err != nil {
		return err
	}

	w.ftch <- ft
	return nil
}

// getReadyWorker .
func (p *Pool) getReadyWorker() (w *worker, err error) {
	w = p.ready.pop()
	if w == nil {
		p.lock.Lock()
		workerID := p.curWorkers
		if p.curWorkers >= p.conf.MaxWorkers {
			err = errors.New("workerpool is full")
			p.lock.Unlock()
			return
		}
		p.curWorkers++
		p.lock.Unlock()

		w = newWorker(workerID)
		go func(w *worker) {
			for {
				ft, ok := <-w.ftch
				if !ok {
					return
				}
				ft.out <- ft.T.Run()
				p.release(w)
			}
		}(w)
	}
	return
}

// close worker
func (p *Pool) close(w *worker) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.curWorkers > 0 {
		p.curWorkers--
	}
	close(w.ftch)
}

// release  worker
func (p *Pool) release(w *worker) {
	if p.state > stateRunning {
		p.close(w)
		return
	}
	w.lastUseTime = time.Now()
	if err := p.ready.push(w); err != nil {
		p.close(w)
	}
}

// clean: clean idle goroutine
func (p *Pool) clean() {
	for {
		size := p.ready.size()
		if size <= p.conf.MinIdleWorkers {
			return
		}

		w := p.ready.pop()
		if w == nil {
			return
		}

		currentTime := time.Now()
		if currentTime.Sub(w.lastUseTime) < p.conf.KeepAlive {
			p.release(w)
			return
		}
		p.close(w)
	}
}

// cleanAll
func (p *Pool) cleanAll() {
	for {
		w := p.ready.pop()
		if w == nil {
			return
		}
		p.release(w)
	}
}
