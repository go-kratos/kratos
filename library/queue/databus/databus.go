package databus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"go-common/library/cache/redis"
	"go-common/library/conf/env"
	"go-common/library/container/pool"
	"go-common/library/log"
	"go-common/library/naming"
	"go-common/library/naming/discovery"
	"go-common/library/net/netutil"
	"go-common/library/net/trace"
	"go-common/library/stat/prom"
	xtime "go-common/library/time"
)

const (
	_appid = "middleware.databus"
)

type dial func() (redis.Conn, error)

// Config databus config.
type Config struct {
	Key          string
	Secret       string
	Group        string
	Topic        string
	Action       string // shoule be "pub" or "sub" or "pubsub"
	Buffer       int
	Name         string // redis name, for trace
	Proto        string
	Addr         string
	Auth         string
	Active       int // pool
	Idle         int // pool
	DialTimeout  xtime.Duration
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
	IdleTimeout  xtime.Duration
	Direct       bool
}

const (
	_family     = "databus"
	_actionSub  = "sub"
	_actionPub  = "pub"
	_actionAll  = "pubsub"
	_cmdPub     = "SET"
	_cmdSub     = "MGET"
	_authFormat = "%s:%s@%s/topic=%s&role=%s"
	_open       = int32(0)
	_closed     = int32(1)
	_scheme     = "databus"
)

var (
	// ErrAction action error.
	ErrAction = errors.New("action unknown")
	// ErrFull chan full
	ErrFull = errors.New("chan full")
	// ErrNoInstance no instances
	ErrNoInstance = errors.New("no databus instances found")
	bk            = netutil.DefaultBackoffConfig
	stats         = prom.LibClient
)

// Message Data.
type Message struct {
	Key       string          `json:"key"`
	Value     json.RawMessage `json:"value"`
	Topic     string          `json:"topic"`
	Partition int32           `json:"partition"`
	Offset    int64           `json:"offset"`
	Timestamp int64           `json:"timestamp"`
	d         *Databus
}

// Commit ack message.
func (m *Message) Commit() (err error) {
	m.d.lock.Lock()
	if m.Offset >= m.d.marked[m.Partition] {
		m.d.marked[m.Partition] = m.Offset
	}
	m.d.lock.Unlock()
	return nil
}

// Databus databus struct.
type Databus struct {
	conf *Config

	d   dial
	p   *redis.Pool
	dis naming.Resolver

	msgs   chan *Message
	lock   sync.RWMutex
	marked map[int32]int64
	idx    int64

	closed int32
}

// New new a databus.
func New(c *Config) *Databus {
	if c.Buffer == 0 {
		c.Buffer = 1024
	}
	d := &Databus{
		conf:   c,
		msgs:   make(chan *Message, c.Buffer),
		marked: make(map[int32]int64),
		closed: _open,
	}

	if !c.Direct && env.DeployEnv != "" && env.DeployEnv != env.DeployEnvDev {
		d.dis = discovery.Build(_appid)
		e := d.dis.Watch()
		select {
		case <-e:
			d.disc()
		case <-time.After(10 * time.Second):
			panic("init discovery err")
		}
		go d.discoveryproc(e)
		log.Info("init databus discvoery info successfully")
	}
	if c.Action == _actionSub || c.Action == _actionAll {
		if d.dis == nil {
			d.d = d.dial
		} else {
			d.d = d.dialInstance
		}
		go d.subproc()
	}
	if c.Action == _actionPub || c.Action == _actionAll {
		// new pool
		d.p = d.redisPool(c)
		if d.dis != nil {
			d.p.New = func(ctx context.Context) (io.Closer, error) {
				return d.dialInstance()
			}
		}
	}
	return d
}

func (d *Databus) redisPool(c *Config) *redis.Pool {
	config := &redis.Config{
		Name:         c.Name,
		Proto:        c.Proto,
		Addr:         c.Addr,
		Auth:         fmt.Sprintf(_authFormat, c.Key, c.Secret, c.Group, c.Topic, c.Action),
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
	}
	config.Config = &pool.Config{
		Active:      c.Active,
		Idle:        c.Idle,
		IdleTimeout: c.IdleTimeout,
	}
	stat := redis.DialStats(statfunc)
	return redis.NewPool(config, stat)
}

func statfunc(cmd string, err *error) func() {
	now := time.Now()
	return func() {
		stats.Timing(fmt.Sprintf("databus:%s", cmd), int64(time.Since(now)/time.Millisecond))
		if err != nil && *err != nil {
			stats.Incr("databus", (*err).Error())
		}
	}
}

func (d *Databus) redisOptions() []redis.DialOption {
	cnop := redis.DialConnectTimeout(time.Duration(d.conf.DialTimeout))
	rdop := redis.DialReadTimeout(time.Duration(d.conf.ReadTimeout))
	wrop := redis.DialWriteTimeout(time.Duration(d.conf.WriteTimeout))
	auop := redis.DialPassword(fmt.Sprintf(_authFormat, d.conf.Key, d.conf.Secret, d.conf.Group, d.conf.Topic, d.conf.Action))
	stat := redis.DialStats(statfunc)
	return []redis.DialOption{cnop, rdop, wrop, auop, stat}
}

func (d *Databus) dial() (redis.Conn, error) {
	return redis.Dial(d.conf.Proto, d.conf.Addr, d.redisOptions()...)
}

func (d *Databus) dialInstance() (redis.Conn, error) {
	if insMap, ok := d.dis.Fetch(context.Background()); ok {
		ins, ok := insMap[env.Zone]
		if !ok || len(ins) == 0 {
			for _, is := range insMap {
				ins = append(ins, is...)
			}
		}
		if len(ins) > 0 {
			var in *naming.Instance
			if d.conf.Action == "pub" {
				i := atomic.AddInt64(&d.idx, 1)
				in = ins[i%int64(len(ins))]
			} else {
				in = ins[rand.Intn(len(ins))]
			}
			for _, addr := range in.Addrs {
				u, err := url.Parse(addr)
				if err == nil && u.Scheme == _scheme {
					return redis.Dial("tcp", u.Host, d.redisOptions()...)
				}
			}
		}
	}
	if d.conf.Proto != "" && d.conf.Addr != "" {
		log.Warn("Databus: no instances(%s,%s) found in discovery,Use config(%s,%s)", _appid, env.Zone, d.conf.Proto, d.conf.Addr)
		return redis.Dial(d.conf.Proto, d.conf.Addr, d.redisOptions()...)
	}
	return nil, ErrNoInstance
}

func (d *Databus) disc() {
	if d.p != nil {
		op := d.p
		np := d.redisPool(d.conf)
		np.New = func(ctx context.Context) (io.Closer, error) {
			return d.dialInstance()
		}
		d.p = np
		op.Close()
		op = nil
		log.Info("discovery event renew redis pool group(%s) topic(%s)", d.conf.Group, d.conf.Topic)
	}
	if insMap, ok := d.dis.Fetch(context.Background()); ok {
		if ins, ok := insMap[env.Zone]; ok && len(ins) > 0 {
			log.Info("get databus instances len(%d)", len(ins))
		}
	}
}

func (d *Databus) discoveryproc(e <-chan struct{}) {
	if d.dis == nil {
		return
	}
	for {
		<-e
		d.disc()
	}
}

func (d *Databus) subproc() {
	var (
		err      error
		r        []byte
		res      [][]byte
		c        redis.Conn
		retry    int
		commited = make(map[int32]int64)
		commit   = make(map[int32]int64)
	)
	for {
		if atomic.LoadInt32(&d.closed) == _closed {
			if c != nil {
				c.Close()
			}
			close(d.msgs)
			return
		}
		if err != nil {
			time.Sleep(bk.Backoff(retry))
			retry++
		} else {
			retry = 0
		}
		if c == nil || c.Err() != nil {
			if c, err = d.d(); err != nil {
				log.Error("redis.Dial(%s@%s) group(%s) retry error(%v)", d.conf.Proto, d.conf.Addr, d.conf.Group, err)
				continue
			}
		}
		d.lock.RLock()
		for k, v := range d.marked {
			if commited[k] != v {
				commit[k] = v
			}
		}
		d.lock.RUnlock()
		// TODO pipeline commit offset
		for k, v := range commit {
			if _, err = c.Do("SET", k, v); err != nil {
				c.Close()
				log.Error("group(%s) conn.Do(SET,%d,%d) commit error(%v)", d.conf.Group, k, v, err)
				break
			}
			delete(commit, k)
			commited[k] = v
		}
		if err != nil {
			continue
		}
		// pull messages
		if res, err = redis.ByteSlices(c.Do(_cmdSub, "")); err != nil {
			c.Close()
			log.Error("group(%s) conn.Do(MGET) error(%v)", d.conf.Group, err)
			continue
		}
		for _, r = range res {
			msg := &Message{d: d}
			if err = json.Unmarshal(r, msg); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", r, err)
				continue
			}
			d.msgs <- msg
		}
	}
}

// Messages get message chan.
func (d *Databus) Messages() <-chan *Message {
	return d.msgs
}

// Send send message to databus.
func (d *Databus) Send(c context.Context, k string, v interface{}) (err error) {
	var b []byte
	// trace info
	if t, ok := trace.FromContext(c); ok {
		t = t.Fork(_family, _cmdPub)
		t.SetTag(trace.String(trace.TagAddress, d.conf.Addr), trace.String(trace.TagComment, k))
		defer t.Finish(&err)
	}
	// send message
	if b, err = json.Marshal(v); err != nil {
		log.Error("json.Marshal(%v) error(%v)", v, err)
		return
	}
	conn := d.p.Get(context.TODO())
	if _, err = conn.Do(_cmdPub, k, b); err != nil {
		log.Error("conn.Do(%s,%s,%s) error(%v)", _cmdPub, k, b, err)
	}
	conn.Close()
	return
}

// Close close databus conn.
func (d *Databus) Close() (err error) {
	if !atomic.CompareAndSwapInt32(&d.closed, _open, _closed) {
		return
	}
	if d.p != nil {
		d.p.Close()
	}
	return nil
}
