package time

import (
	"sync"
	itime "time"

	"go-common/library/log"
)

const (
	timerFormat      = "2006-01-02 15:04:05"
	infiniteDuration = itime.Duration(1<<63 - 1)
)

var (
	timerLazyDelay = 300 * itime.Millisecond
)

// TimerData timer data.
type TimerData struct {
	Key    string
	expire itime.Time
	fn     func()
	index  int
	next   *TimerData
}

// Delay delay duration.
func (td *TimerData) Delay() itime.Duration {
	return td.expire.Sub(itime.Now())
}

// ExpireString expire string.
func (td *TimerData) ExpireString() string {
	return td.expire.Format(timerFormat)
}

// Timer timer.
type Timer struct {
	lock   sync.Mutex
	free   *TimerData
	timers []*TimerData
	signal *itime.Timer
	num    int
}

// NewTimer new a timer.
// A heap must be initialized before any of the heap operations
// can be used. Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// Its complexity is O(n) where n = h.Len().
//
func NewTimer(num int) (t *Timer) {
	t = new(Timer)
	t.init(num)
	return t
}

// Init init the timer.
func (t *Timer) Init(num int) {
	t.init(num)
}

func (t *Timer) init(num int) {
	t.signal = itime.NewTimer(infiniteDuration)
	t.timers = make([]*TimerData, 0, num)
	t.num = num
	t.grow()
	go t.start()
}

func (t *Timer) grow() {
	var (
		i   int
		td  *TimerData
		tds = make([]TimerData, t.num)
	)
	t.free = &(tds[0])
	td = t.free
	for i = 1; i < t.num; i++ {
		td.next = &(tds[i])
		td = td.next
	}
	td.next = nil
}

// get get a free timer data.
func (t *Timer) get() (td *TimerData) {
	if td = t.free; td == nil {
		t.grow()
		td = t.free
	}
	t.free = td.next
	return
}

// put put back a timer data.
func (t *Timer) put(td *TimerData) {
	td.fn = nil
	td.next = t.free
	t.free = td
}

// Add add the element x onto the heap. The complexity is
// O(log(n)) where n = h.Len().
func (t *Timer) Add(expire itime.Duration, fn func()) (td *TimerData) {
	t.lock.Lock()
	td = t.get()
	td.expire = itime.Now().Add(expire)
	td.fn = fn
	t.add(td)
	t.lock.Unlock()
	return
}

// Del removes the element at index i from the heap.
// The complexity is O(log(n)) where n = h.Len().
func (t *Timer) Del(td *TimerData) {
	t.lock.Lock()
	t.del(td)
	t.put(td)
	t.lock.Unlock()
}

// Push pushes the element x onto the heap. The complexity is
// O(log(n)) where n = h.Len().
func (t *Timer) add(td *TimerData) {
	var d itime.Duration
	td.index = len(t.timers)
	// add to the minheap last node
	t.timers = append(t.timers, td)
	t.up(td.index)
	if td.index == 0 {
		// if first node, signal start goroutine
		d = td.Delay()
		t.signal.Reset(d)
		if Debug {
			log.Info("timer: add reset delay %d ms", int64(d)/int64(itime.Millisecond))
		}
	}
	if Debug {
		log.Info("timer: push item key: %s, expire: %s, index: %d", td.Key, td.ExpireString(), td.index)
	}
}

func (t *Timer) del(td *TimerData) {
	var (
		i    = td.index
		last = len(t.timers) - 1
	)
	if i < 0 || i > last || t.timers[i] != td {
		// already remove, usually by expire
		if Debug {
			log.Info("timer del i: %d, last: %d, %p", i, last, td)
		}
		return
	}
	if i != last {
		t.swap(i, last)
		t.down(i, last)
		t.up(i)
	}
	// remove item is the last node
	t.timers[last].index = -1 // for safety
	t.timers = t.timers[:last]
	if Debug {
		log.Info("timer: remove item key: %s, expire: %s, index: %d", td.Key, td.ExpireString(), td.index)
	}
}

// Set update timer data.
func (t *Timer) Set(td *TimerData, expire itime.Duration) {
	t.lock.Lock()
	t.del(td)
	td.expire = itime.Now().Add(expire)
	t.add(td)
	t.lock.Unlock()
}

// start start the timer.
func (t *Timer) start() {
	for {
		t.expire()
		<-t.signal.C
	}
}

// expire removes the minimum element (according to Less) from the heap.
// The complexity is O(log(n)) where n = max.
// It is equivalent to Del(0).
func (t *Timer) expire() {
	var (
		fn func()
		td *TimerData
		d  itime.Duration
	)
	t.lock.Lock()
	for {
		if len(t.timers) == 0 {
			d = infiniteDuration
			if Debug {
				log.Info("timer: no other instance")
			}
			break
		}
		td = t.timers[0]
		if d = td.Delay(); d > 0 {
			break
		}
		fn = td.fn
		// let caller put back
		t.del(td)
		t.lock.Unlock()
		if fn == nil {
			log.Warn("expire timer no fn")
		} else {
			if Debug {
				log.Info("timer key: %s, expire: %s, index: %d expired, call fn", td.Key, td.ExpireString(), td.index)
			}
			fn()
		}
		t.lock.Lock()
	}
	t.signal.Reset(d)
	if Debug {
		log.Info("timer: expier reset delay %d ms", int64(d)/int64(itime.Millisecond))
	}
	t.lock.Unlock()
}

func (t *Timer) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i <= j || !t.less(j, i) {
			break
		}
		t.swap(i, j)
		j = i
	}
}

func (t *Timer) down(i, n int) {
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && !t.less(j1, j2) {
			j = j2 // = 2*i + 2  // right child
		}
		if !t.less(j, i) {
			break
		}
		t.swap(i, j)
		i = j
	}
}

func (t *Timer) less(i, j int) bool {
	return t.timers[i].expire.Before(t.timers[j].expire)
}

func (t *Timer) swap(i, j int) {
	t.timers[i], t.timers[j] = t.timers[j], t.timers[i]
	t.timers[i].index = i
	t.timers[j].index = j
}
