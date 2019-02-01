package breaker

import (
	"sync"
	"time"

	xtime "go-common/library/time"
)

// Config broker config.
type Config struct {
	SwitchOff bool // breaker switch,default off.

	// Hystrix
	Ratio float32
	Sleep xtime.Duration

	// Google
	K float64

	Window  xtime.Duration
	Bucket  int
	Request int64
}

func (conf *Config) fix() {
	if conf.K == 0 {
		conf.K = 1.5
	}
	if conf.Request == 0 {
		conf.Request = 100
	}
	if conf.Ratio == 0 {
		conf.Ratio = 0.5
	}
	if conf.Sleep == 0 {
		conf.Sleep = xtime.Duration(500 * time.Millisecond)
	}
	if conf.Bucket == 0 {
		conf.Bucket = 10
	}
	if conf.Window == 0 {
		conf.Window = xtime.Duration(3 * time.Second)
	}
}

// Breaker is a CircuitBreaker pattern.
// FIXME on int32 atomic.LoadInt32(&b.on) == _switchOn
type Breaker interface {
	Allow() error
	MarkSuccess()
	MarkFailed()
}

// Group represents a class of CircuitBreaker and forms a namespace in which
// units of CircuitBreaker.
type Group struct {
	mu   sync.RWMutex
	brks map[string]Breaker
	conf *Config
}

const (
	// StateOpen when circuit breaker open, request not allowed, after sleep
	// some duration, allow one single request for testing the health, if ok
	// then state reset to closed, if not continue the step.
	StateOpen int32 = iota
	// StateClosed when circuit breaker closed, request allowed, the breaker
	// calc the succeed ratio, if request num greater request setting and
	// ratio lower than the setting ratio, then reset state to open.
	StateClosed
	// StateHalfopen when circuit breaker open, after slepp some duration, allow
	// one request, but not state closed.
	StateHalfopen

	//_switchOn int32 = iota
	// _switchOff
)

var (
	_mu   sync.RWMutex
	_conf = &Config{
		Window:  xtime.Duration(3 * time.Second),
		Bucket:  10,
		Request: 100,

		Sleep: xtime.Duration(500 * time.Millisecond),
		Ratio: 0.5,
		// Percentage of failures must be lower than 33.33%
		K: 1.5,

		// Pattern: "",
	}
	_group = NewGroup(_conf)
)

// Init init global breaker config, also can reload config after first time call.
func Init(conf *Config) {
	if conf == nil {
		return
	}
	_mu.Lock()
	_conf = conf
	_mu.Unlock()
}

// Go runs your function while tracking the breaker state of default group.
func Go(name string, run, fallback func() error) error {
	breaker := _group.Get(name)
	if err := breaker.Allow(); err != nil {
		return fallback()
	}
	return run()
}

// newBreaker new a breaker.
func newBreaker(c *Config) (b Breaker) {
	// factory
	return newSRE(c)
}

// NewGroup new a breaker group container, if conf nil use default conf.
func NewGroup(conf *Config) *Group {
	if conf == nil {
		_mu.RLock()
		conf = _conf
		_mu.RUnlock()
	} else {
		conf.fix()
	}
	return &Group{
		conf: conf,
		brks: make(map[string]Breaker),
	}
}

// Get get a breaker by a specified key, if breaker not exists then make a new one.
func (g *Group) Get(key string) Breaker {
	g.mu.RLock()
	brk, ok := g.brks[key]
	conf := g.conf
	g.mu.RUnlock()
	if ok {
		return brk
	}
	// NOTE here may new multi breaker for rarely case, let gc drop it.
	brk = newBreaker(conf)
	g.mu.Lock()
	if _, ok = g.brks[key]; !ok {
		g.brks[key] = brk
	}
	g.mu.Unlock()
	return brk
}

// Reload reload the group by specified config, this may let all inner breaker
// reset to a new one.
func (g *Group) Reload(conf *Config) {
	if conf == nil {
		return
	}
	conf.fix()
	g.mu.Lock()
	g.conf = conf
	g.brks = make(map[string]Breaker, len(g.brks))
	g.mu.Unlock()
}

// Go runs your function while tracking the breaker state of group.
func (g *Group) Go(name string, run, fallback func() error) error {
	breaker := g.Get(name)
	if err := breaker.Allow(); err != nil {
		return fallback()
	}
	return run()
}
