package pipeline

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/bilibili/kratos/pkg/net/metadata"
	"github.com/bilibili/kratos/pkg/stat/metric"
	xtime "github.com/bilibili/kratos/pkg/time"
)

// ErrFull channel full error
var ErrFull = errors.New("channel full")

const _metricNamespace = "sync"
const _metricSubSystem = "pipeline"

var (
	_metricCount = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: _metricNamespace,
		Subsystem: _metricSubSystem,
		Name:      "process_count",
		Help:      "process count",
		Labels:    []string{"name", "chan"},
	})
	_metricChanLen = metric.NewGaugeVec(&metric.GaugeVecOpts{
		Namespace: _metricNamespace,
		Subsystem: _metricSubSystem,
		Name:      "chan_len",
		Help:      "channel length",
		Labels:    []string{"name", "chan"},
	})
)

type message struct {
	key   string
	value interface{}
}

// Pipeline pipeline struct
type Pipeline struct {
	Do          func(c context.Context, index int, values map[string][]interface{})
	Split       func(key string) int
	chans       []chan *message
	mirrorChans []chan *message
	config      *Config
	wait        sync.WaitGroup
	name        string
}

// Config Pipeline config
type Config struct {
	// MaxSize merge size
	MaxSize int
	// Interval merge interval
	Interval xtime.Duration
	// Buffer channel size
	Buffer int
	// Worker channel number
	Worker int
	// Name use for metrics
	Name string
}

func (c *Config) fix() {
	if c.MaxSize <= 0 {
		c.MaxSize = 1000
	}
	if c.Interval <= 0 {
		c.Interval = xtime.Duration(time.Second)
	}
	if c.Buffer <= 0 {
		c.Buffer = 1000
	}
	if c.Worker <= 0 {
		c.Worker = 10
	}
	if c.Name == "" {
		c.Name = "anonymous"
	}
}

// NewPipeline new pipline
func NewPipeline(config *Config) (res *Pipeline) {
	if config == nil {
		config = &Config{}
	}
	config.fix()
	res = &Pipeline{
		chans:       make([]chan *message, config.Worker),
		mirrorChans: make([]chan *message, config.Worker),
		config:      config,
		name:        config.Name,
	}
	for i := 0; i < config.Worker; i++ {
		res.chans[i] = make(chan *message, config.Buffer)
		res.mirrorChans[i] = make(chan *message, config.Buffer)
	}
	return
}

// Start start all mergeproc
func (p *Pipeline) Start() {
	if p.Do == nil {
		panic("pipeline: do func is nil")
	}
	if p.Split == nil {
		panic("pipeline: split func is nil")
	}
	var mirror bool
	p.wait.Add(len(p.chans) + len(p.mirrorChans))
	for i, ch := range p.chans {
		go p.mergeproc(mirror, i, ch)
	}
	mirror = true
	for i, ch := range p.mirrorChans {
		go p.mergeproc(mirror, i, ch)
	}
}

// SyncAdd sync add a value to channal, channel shard in split method
func (p *Pipeline) SyncAdd(c context.Context, key string, value interface{}) (err error) {
	ch, msg := p.add(c, key, value)
	select {
	case ch <- msg:
	case <-c.Done():
		err = c.Err()
	}
	return
}

// Add async add a value to channal, channel shard in split method
func (p *Pipeline) Add(c context.Context, key string, value interface{}) (err error) {
	ch, msg := p.add(c, key, value)
	select {
	case ch <- msg:
	default:
		err = ErrFull
	}
	return
}

func (p *Pipeline) add(c context.Context, key string, value interface{}) (ch chan *message, m *message) {
	shard := p.Split(key) % p.config.Worker
	if metadata.String(c, metadata.Mirror) != "" {
		ch = p.mirrorChans[shard]
	} else {
		ch = p.chans[shard]
	}
	m = &message{key: key, value: value}
	return
}

// Close all goroutinue
func (p *Pipeline) Close() (err error) {
	for _, ch := range p.chans {
		ch <- nil
	}
	for _, ch := range p.mirrorChans {
		ch <- nil
	}
	p.wait.Wait()
	return
}

func (p *Pipeline) mergeproc(mirror bool, index int, ch <-chan *message) {
	defer p.wait.Done()
	var (
		m         *message
		vals      = make(map[string][]interface{}, p.config.MaxSize)
		closed    bool
		count     int
		inteval   = p.config.Interval
		oldTicker = true
	)
	if index > 0 {
		inteval = xtime.Duration(int64(index) * (int64(p.config.Interval) / int64(p.config.Worker)))
	}
	ticker := time.NewTicker(time.Duration(inteval))
	for {
		select {
		case m = <-ch:
			if m == nil {
				closed = true
				break
			}
			count++
			vals[m.key] = append(vals[m.key], m.value)
			if count >= p.config.MaxSize {
				break
			}
			continue
		case <-ticker.C:
			if oldTicker {
				ticker.Stop()
				ticker = time.NewTicker(time.Duration(p.config.Interval))
				oldTicker = false
			}
		}
		name := p.name
		process := count
		if len(vals) > 0 {
			ctx := context.Background()
			if mirror {
				ctx = metadata.NewContext(ctx, metadata.MD{metadata.Mirror: "1"})
				name = "mirror_" + name
			}
			p.Do(ctx, index, vals)
			vals = make(map[string][]interface{}, p.config.MaxSize)
			count = 0
		}
		_metricChanLen.Set(float64(len(ch)), name, strconv.Itoa(index))
		_metricCount.Add(float64(process), name, strconv.Itoa(index))
		if closed {
			ticker.Stop()
			return
		}
	}
}
