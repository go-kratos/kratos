// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package redis

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/go-kratos/kratos/pkg/container/pool"
	"github.com/go-kratos/kratos/pkg/net/trace"
	xtime "github.com/go-kratos/kratos/pkg/time"
)

var beginTime, _ = time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")

var (
	errConnClosed = errors.New("redigo: connection closed")
)

// Pool .
type Pool struct {
	*pool.Slice
	// config
	c *Config
	// statfunc
	statfunc func(name, addr, cmd string, t time.Time, err error) func()
}

// NewPool creates a new pool.
func NewPool(c *Config, options ...DialOption) (p *Pool) {
	if c.DialTimeout <= 0 || c.ReadTimeout <= 0 || c.WriteTimeout <= 0 {
		panic("must config redis timeout")
	}
	if c.SlowLog <= 0 {
		c.SlowLog = xtime.Duration(250 * time.Millisecond)
	}
	ops := []DialOption{
		DialConnectTimeout(time.Duration(c.DialTimeout)),
		DialReadTimeout(time.Duration(c.ReadTimeout)),
		DialWriteTimeout(time.Duration(c.WriteTimeout)),
		DialPassword(c.Auth),
		DialDatabase(c.Db),
	}
	ops = append(ops, options...)
	p1 := pool.NewSlice(c.Config)

	// new pool
	p1.New = func(ctx context.Context) (io.Closer, error) {
		conn, err := Dial(c.Proto, c.Addr, ops...)
		if err != nil {
			return nil, err
		}
		return &traceConn{
			Conn:             conn,
			connTags:         []trace.Tag{trace.TagString(trace.TagPeerAddress, c.Addr)},
			slowLogThreshold: time.Duration(c.SlowLog),
		}, nil
	}
	p = &Pool{Slice: p1, c: c, statfunc: pstat}
	return
}

// Get gets a connection. The application must close the returned connection.
// This method always returns a valid connection so that applications can defer
// error handling to the first use of the connection. If there is an error
// getting an underlying connection, then the connection Err, Do, Send, Flush
// and Receive methods return that error.
func (p *Pool) Get(ctx context.Context) Conn {
	c, err := p.Slice.Get(ctx)
	if err != nil {
		return errorConnection{err}
	}
	c1, _ := c.(Conn)
	return &pooledConnection{p: p, c: c1.WithContext(ctx), rc: c1, now: beginTime}
}

// Close releases the resources used by the pool.
func (p *Pool) Close() error {
	return p.Slice.Close()
}

type pooledConnection struct {
	p     *Pool
	rc    Conn
	c     Conn
	state int

	now  time.Time
	cmds []string
}

var (
	sentinel     []byte
	sentinelOnce sync.Once
)

func initSentinel() {
	p := make([]byte, 64)
	if _, err := rand.Read(p); err == nil {
		sentinel = p
	} else {
		h := sha1.New()
		io.WriteString(h, "Oops, rand failed. Use time instead.")
		io.WriteString(h, strconv.FormatInt(time.Now().UnixNano(), 10))
		sentinel = h.Sum(nil)
	}
}

// SetStatFunc set stat func.
func (p *Pool) SetStatFunc(fn func(name, addr, cmd string, t time.Time, err error) func()) {
	p.statfunc = fn
}

func pstat(name, addr, cmd string, t time.Time, err error) func() {
	return func() {
		_metricReqDur.Observe(int64(time.Since(t)/time.Millisecond), name, addr, cmd)
		if err != nil {
			if msg := formatErr(err, name, addr); msg != "" {
				_metricReqErr.Inc(name, addr, cmd, msg)
			}
			return
		}
		_metricHits.Inc(name, addr)
	}
}

func (pc *pooledConnection) Close() error {
	c := pc.c
	if _, ok := c.(errorConnection); ok {
		return nil
	}
	pc.c = errorConnection{errConnClosed}

	if pc.state&MultiState != 0 {
		c.Send("DISCARD")
		pc.state &^= (MultiState | WatchState)
	} else if pc.state&WatchState != 0 {
		c.Send("UNWATCH")
		pc.state &^= WatchState
	}
	if pc.state&SubscribeState != 0 {
		c.Send("UNSUBSCRIBE")
		c.Send("PUNSUBSCRIBE")
		// To detect the end of the message stream, ask the server to echo
		// a sentinel value and read until we see that value.
		sentinelOnce.Do(initSentinel)
		c.Send("ECHO", sentinel)
		c.Flush()
		for {
			p, err := c.Receive()
			if err != nil {
				break
			}
			if p, ok := p.([]byte); ok && bytes.Equal(p, sentinel) {
				pc.state &^= SubscribeState
				break
			}
		}
	}
	_, err := c.Do("")
	pc.p.Slice.Put(context.Background(), pc.rc, pc.state != 0 || c.Err() != nil)
	return err
}

func (pc *pooledConnection) Err() error {
	return pc.c.Err()
}

func (pc *pooledConnection) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	now := time.Now()
	ci := LookupCommandInfo(commandName)
	pc.state = (pc.state | ci.Set) &^ ci.Clear
	reply, err = pc.c.Do(commandName, args...)
	if pc.p.statfunc != nil {
		pc.p.statfunc(pc.p.c.Name, pc.p.c.Addr, commandName, now, err)()
	}
	return
}

func (pc *pooledConnection) Send(commandName string, args ...interface{}) (err error) {
	ci := LookupCommandInfo(commandName)
	pc.state = (pc.state | ci.Set) &^ ci.Clear
	if pc.now.Equal(beginTime) {
		// mark first send time
		pc.now = time.Now()
	}
	pc.cmds = append(pc.cmds, commandName)
	return pc.c.Send(commandName, args...)
}

func (pc *pooledConnection) Flush() error {
	return pc.c.Flush()
}

func (pc *pooledConnection) Receive() (reply interface{}, err error) {
	reply, err = pc.c.Receive()
	if len(pc.cmds) > 0 {
		cmd := pc.cmds[0]
		pc.cmds = pc.cmds[1:]
		if pc.p.statfunc != nil {
			pc.p.statfunc(pc.p.c.Name, pc.p.c.Addr, cmd, pc.now, err)()
		}
	}
	return
}

func (pc *pooledConnection) WithContext(ctx context.Context) Conn {
	return pc
}

type errorConnection struct{ err error }

func (ec errorConnection) Do(string, ...interface{}) (interface{}, error) {
	return nil, ec.err
}
func (ec errorConnection) Send(string, ...interface{}) error { return ec.err }
func (ec errorConnection) Err() error                        { return ec.err }
func (ec errorConnection) Close() error                      { return ec.err }
func (ec errorConnection) Flush() error                      { return ec.err }
func (ec errorConnection) Receive() (interface{}, error)     { return nil, ec.err }
func (ec errorConnection) WithContext(context.Context) Conn  { return ec }
