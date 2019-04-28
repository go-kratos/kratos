package databusutil

import (
	"runtime"
	"sync"
	"time"

	"go-common/library/queue/databus"
	xtime "go-common/library/time"
)

const (
	_stateStarted = 1
	_stateClosed  = 2
)

// Config the config is the base configuration for initiating a new group.
type Config struct {
	// Size merge size
	Size int
	// Num merge goroutine num
	Num int
	// Ticker duration of submit merges when no new message
	Ticker xtime.Duration
	// Chan size of merge chan and done chan
	Chan int
}

func (c *Config) fix() {
	if c.Size <= 0 {
		c.Size = 1024
	}
	if int64(c.Ticker) <= 0 {
		c.Ticker = xtime.Duration(time.Second * 5)
	}
	if c.Num <= 0 {
		c.Num = runtime.GOMAXPROCS(0)
	}
	if c.Chan <= 0 {
		c.Chan = 1024
	}
}

type message struct {
	next   *message
	data   *databus.Message
	object interface{}
	done   bool
}

// Group group.
type Group struct {
	c          *Config
	head, last *message
	state      int
	mu         sync.Mutex

	mc  []chan *message // merge chan
	dc  chan []*message // done chan
	qc  chan struct{}   // quit chan
	msg <-chan *databus.Message

	New   func(msg *databus.Message) (interface{}, error)
	Split func(msg *databus.Message, data interface{}) int
	Do    func(msgs []interface{})

	pool *sync.Pool
}

// NewGroup new a group.
func NewGroup(c *Config, m <-chan *databus.Message) *Group {
	// NOTE if c || m == nil runtime panic
	if c == nil {
		c = new(Config)
	}
	c.fix()
	g := &Group{
		c:   c,
		msg: m,

		mc: make([]chan *message, c.Num),
		dc: make(chan []*message, c.Chan),
		qc: make(chan struct{}),

		pool: &sync.Pool{
			New: func() interface{} {
				return new(message)
			},
		},
	}
	for i := 0; i < c.Num; i++ {
		g.mc[i] = make(chan *message, c.Chan)
	}
	return g
}

// Start start group, it is safe for concurrent use by multiple goroutines.
func (g *Group) Start() {
	g.mu.Lock()
	if g.state == _stateStarted {
		g.mu.Unlock()
		return
	}
	g.state = _stateStarted
	g.mu.Unlock()
	go g.consumeproc()
	for i := 0; i < g.c.Num; i++ {
		go g.mergeproc(g.mc[i])
	}
	go g.commitproc()
}

// Close close group, it is safe for concurrent use by multiple goroutines.
func (g *Group) Close() (err error) {
	g.mu.Lock()
	if g.state == _stateClosed {
		g.mu.Unlock()
		return
	}
	g.state = _stateClosed
	g.mu.Unlock()
	close(g.qc)
	return
}

func (g *Group) message() *message {
	return g.pool.Get().(*message)
}

func (g *Group) freeMessage(m *message) {
	*m = message{}
	g.pool.Put(m)
}

func (g *Group) consumeproc() {
	var (
		ok  bool
		err error
		msg *databus.Message
	)
	for {
		select {
		case <-g.qc:
			return
		case msg, ok = <-g.msg:
			if !ok {
				g.Close()
				return
			}
		}
		// marked head to first commit
		m := g.message()
		m.data = msg
		if m.object, err = g.New(msg); err != nil {
			g.freeMessage(m)
			continue
		}
		g.mu.Lock()
		if g.head == nil {
			g.head = m
			g.last = m
		} else {
			g.last.next = m
			g.last = m
		}
		g.mu.Unlock()
		g.mc[g.Split(m.data, m.object)%g.c.Num] <- m
	}
}

func (g *Group) mergeproc(mc <-chan *message) {
	ticker := time.NewTicker(time.Duration(g.c.Ticker))
	msgs := make([]interface{}, 0, g.c.Size)
	marks := make([]*message, 0, g.c.Size)
	for {
		select {
		case <-g.qc:
			return
		case msg := <-mc:
			msgs = append(msgs, msg.object)
			marks = append(marks, msg)
			if len(msgs) < g.c.Size {
				continue
			}
		case <-ticker.C:
		}
		if len(msgs) > 0 {
			g.Do(msgs)
			msgs = make([]interface{}, 0, g.c.Size)
		}
		if len(marks) > 0 {
			g.dc <- marks
			marks = make([]*message, 0, g.c.Size)
		}
	}
}

func (g *Group) commitproc() {
	commits := make(map[int32]*databus.Message)
	for {
		select {
		case <-g.qc:
			return
		case done := <-g.dc:
			// merge partitions to commit offset
			for _, d := range done {
				d.done = true
			}
			g.mu.Lock()
			for g.head != nil && g.head.done {
				cur := g.head
				commits[cur.data.Partition] = cur.data
				g.head = cur.next
				g.freeMessage(cur)
			}
			g.mu.Unlock()
			for k, m := range commits {
				m.Commit()
				delete(commits, k)
			}
		}
	}
}
