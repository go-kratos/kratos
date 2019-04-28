// Copyright 2011 Gary Burd
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
	"context"
	"errors"
	"io"
	"reflect"
	"sync"
	"testing"
	"time"

	"go-common/library/container/pool"
)

type poolTestConn struct {
	d   *poolDialer
	err error
	Conn
}

func (c *poolTestConn) Close() error {
	c.d.mu.Lock()
	c.d.open--
	c.d.mu.Unlock()
	return c.Conn.Close()
}

func (c *poolTestConn) Err() error { return c.err }

func (c *poolTestConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	if commandName == "ERR" {
		c.err = args[0].(error)
		commandName = "PING"
	}
	if commandName != "" {
		c.d.commands = append(c.d.commands, commandName)
	}
	return c.Conn.Do(commandName, args...)
}

func (c *poolTestConn) Send(commandName string, args ...interface{}) error {
	c.d.commands = append(c.d.commands, commandName)
	return c.Conn.Send(commandName, args...)
}

type poolDialer struct {
	mu       sync.Mutex
	t        *testing.T
	dialed   int
	open     int
	commands []string
	dialErr  error
}

func (d *poolDialer) dial() (Conn, error) {
	d.mu.Lock()
	d.dialed++
	dialErr := d.dialErr
	d.mu.Unlock()
	if dialErr != nil {
		return nil, d.dialErr
	}
	c, err := DialDefaultServer()
	if err != nil {
		return nil, err
	}
	d.mu.Lock()
	d.open++
	d.mu.Unlock()
	return &poolTestConn{d: d, Conn: c}, nil
}

func (d *poolDialer) check(message string, p *Pool, dialed, open int) {
	d.mu.Lock()
	if d.dialed != dialed {
		d.t.Errorf("%s: dialed=%d, want %d", message, d.dialed, dialed)
	}
	if d.open != open {
		d.t.Errorf("%s: open=%d, want %d", message, d.open, open)
	}
	//	if active := p.ActiveCount(); active != open {
	//		d.t.Errorf("%s: active=%d, want %d", message, active, open)
	//	}
	d.mu.Unlock()
}

func TestPoolReuse(t *testing.T) {
	d := poolDialer{t: t}
	p := NewPool(config)
	p.Slice.New = func(ctx context.Context) (io.Closer, error) {
		return d.dial()
	}

	for i := 0; i < 10; i++ {
		c1 := p.Get(context.TODO())
		c1.Do("PING")
		c2 := p.Get(context.TODO())
		c2.Do("PING")
		c1.Close()
		c2.Close()
	}

	d.check("before close", p, 2, 2)
	p.Close()
	d.check("after close", p, 2, 0)
}

func TestPoolMaxIdle(t *testing.T) {
	d := poolDialer{t: t}
	p := NewPool(config)
	p.Slice.New = func(ctx context.Context) (io.Closer, error) {
		return d.dial()
	}
	defer p.Close()

	for i := 0; i < 10; i++ {
		c1 := p.Get(context.TODO())
		c1.Do("PING")
		c2 := p.Get(context.TODO())
		c2.Do("PING")
		c3 := p.Get(context.TODO())
		c3.Do("PING")
		c1.Close()
		c2.Close()
		c3.Close()
	}
	d.check("before close", p, 12, 2)
	p.Close()
	d.check("after close", p, 12, 0)
}

func TestPoolError(t *testing.T) {
	d := poolDialer{t: t}
	p := NewPool(config)
	p.Slice.New = func(ctx context.Context) (io.Closer, error) {
		return d.dial()
	}
	defer p.Close()

	c := p.Get(context.TODO())
	c.Do("ERR", io.EOF)
	if c.Err() == nil {
		t.Errorf("expected c.Err() != nil")
	}
	c.Close()

	c = p.Get(context.TODO())
	c.Do("ERR", io.EOF)
	c.Close()

	d.check(".", p, 2, 0)
}

func TestPoolClose(t *testing.T) {
	d := poolDialer{t: t}
	p := NewPool(config)
	p.Slice.New = func(ctx context.Context) (io.Closer, error) {
		return d.dial()
	}
	defer p.Close()

	c1 := p.Get(context.TODO())
	c1.Do("PING")
	c2 := p.Get(context.TODO())
	c2.Do("PING")
	c3 := p.Get(context.TODO())
	c3.Do("PING")

	c1.Close()
	if _, err := c1.Do("PING"); err == nil {
		t.Errorf("expected error after connection closed")
	}

	c2.Close()
	c2.Close()

	p.Close()

	d.check("after pool close", p, 3, 1)

	if _, err := c1.Do("PING"); err == nil {
		t.Errorf("expected error after connection and pool closed")
	}

	c3.Close()

	d.check("after conn close", p, 3, 0)

	c1 = p.Get(context.TODO())
	if _, err := c1.Do("PING"); err == nil {
		t.Errorf("expected error after pool closed")
	}
}

func TestPoolConcurrenSendReceive(t *testing.T) {
	p := NewPool(config)
	p.Slice.New = func(ctx context.Context) (io.Closer, error) {
		return DialDefaultServer()
	}
	defer p.Close()

	c := p.Get(context.TODO())
	done := make(chan error, 1)
	go func() {
		_, err := c.Receive()
		done <- err
	}()
	c.Send("PING")
	c.Flush()
	err := <-done
	if err != nil {
		t.Fatalf("Receive() returned error %v", err)
	}
	_, err = c.Do("")
	if err != nil {
		t.Fatalf("Do() returned error %v", err)
	}
	c.Close()
}

func TestPoolMaxActive(t *testing.T) {
	d := poolDialer{t: t}
	config.Config = &pool.Config{
		Active: 2,
		Idle:   2,
	}
	p := NewPool(config)
	p.Slice.New = func(ctx context.Context) (io.Closer, error) {
		return d.dial()
	}
	defer p.Close()

	c1 := p.Get(context.TODO())
	c1.Do("PING")
	c2 := p.Get(context.TODO())
	c2.Do("PING")

	d.check("1", p, 2, 2)

	c3 := p.Get(context.TODO())
	if _, err := c3.Do("PING"); err != pool.ErrPoolExhausted {
		t.Errorf("expected pool exhausted")
	}

	c3.Close()
	d.check("2", p, 2, 2)
	c2.Close()
	d.check("3", p, 2, 2)

	c3 = p.Get(context.TODO())
	if _, err := c3.Do("PING"); err != nil {
		t.Errorf("expected good channel, err=%v", err)
	}
	c3.Close()

	d.check("4", p, 2, 2)
}

func TestPoolMonitorCleanup(t *testing.T) {
	d := poolDialer{t: t}
	p := NewPool(config)
	p.Slice.New = func(ctx context.Context) (io.Closer, error) {
		return d.dial()
	}
	defer p.Close()

	c := p.Get(context.TODO())
	c.Send("MONITOR")
	c.Close()

	d.check("", p, 1, 0)
}

func TestPoolPubSubCleanup(t *testing.T) {
	d := poolDialer{t: t}
	p := NewPool(config)
	p.Slice.New = func(ctx context.Context) (io.Closer, error) {
		return d.dial()
	}
	defer p.Close()

	c := p.Get(context.TODO())
	c.Send("SUBSCRIBE", "x")
	c.Close()

	want := []string{"SUBSCRIBE", "UNSUBSCRIBE", "PUNSUBSCRIBE", "ECHO"}
	if !reflect.DeepEqual(d.commands, want) {
		t.Errorf("got commands %v, want %v", d.commands, want)
	}
	d.commands = nil

	c = p.Get(context.TODO())
	c.Send("PSUBSCRIBE", "x*")
	c.Close()

	want = []string{"PSUBSCRIBE", "UNSUBSCRIBE", "PUNSUBSCRIBE", "ECHO"}
	if !reflect.DeepEqual(d.commands, want) {
		t.Errorf("got commands %v, want %v", d.commands, want)
	}
	d.commands = nil
}

func TestPoolTransactionCleanup(t *testing.T) {
	d := poolDialer{t: t}
	p := NewPool(config)
	p.Slice.New = func(ctx context.Context) (io.Closer, error) {
		return d.dial()
	}
	defer p.Close()

	c := p.Get(context.TODO())
	c.Do("WATCH", "key")
	c.Do("PING")
	c.Close()

	want := []string{"WATCH", "PING", "UNWATCH"}
	if !reflect.DeepEqual(d.commands, want) {
		t.Errorf("got commands %v, want %v", d.commands, want)
	}
	d.commands = nil

	c = p.Get(context.TODO())
	c.Do("WATCH", "key")
	c.Do("UNWATCH")
	c.Do("PING")
	c.Close()

	want = []string{"WATCH", "UNWATCH", "PING"}
	if !reflect.DeepEqual(d.commands, want) {
		t.Errorf("got commands %v, want %v", d.commands, want)
	}
	d.commands = nil

	c = p.Get(context.TODO())
	c.Do("WATCH", "key")
	c.Do("MULTI")
	c.Do("PING")
	c.Close()

	want = []string{"WATCH", "MULTI", "PING", "DISCARD"}
	if !reflect.DeepEqual(d.commands, want) {
		t.Errorf("got commands %v, want %v", d.commands, want)
	}
	d.commands = nil

	c = p.Get(context.TODO())
	c.Do("WATCH", "key")
	c.Do("MULTI")
	c.Do("DISCARD")
	c.Do("PING")
	c.Close()

	want = []string{"WATCH", "MULTI", "DISCARD", "PING"}
	if !reflect.DeepEqual(d.commands, want) {
		t.Errorf("got commands %v, want %v", d.commands, want)
	}
	d.commands = nil

	c = p.Get(context.TODO())
	c.Do("WATCH", "key")
	c.Do("MULTI")
	c.Do("EXEC")
	c.Do("PING")
	c.Close()

	want = []string{"WATCH", "MULTI", "EXEC", "PING"}
	if !reflect.DeepEqual(d.commands, want) {
		t.Errorf("got commands %v, want %v", d.commands, want)
	}
	d.commands = nil
}

func startGoroutines(p *Pool, cmd string, args ...interface{}) chan error {
	errs := make(chan error, 10)
	for i := 0; i < cap(errs); i++ {
		go func() {
			c := p.Get(context.TODO())
			_, err := c.Do(cmd, args...)
			errs <- err
			c.Close()
		}()
	}

	// Wait for goroutines to block.
	time.Sleep(time.Second / 4)

	return errs
}

func TestWaitPoolDialError(t *testing.T) {
	testErr := errors.New("test")
	d := poolDialer{t: t}
	config1 := getConfig()
	config1.Config = &pool.Config{
		Active: 1,
		Idle:   1,
		Wait:   true,
	}
	p := NewPool(config1)
	p.Slice.New = func(ctx context.Context) (io.Closer, error) {
		return d.dial()
	}
	defer p.Close()

	c := p.Get(context.TODO())
	errs := startGoroutines(p, "ERR", testErr)
	d.check("before close", p, 1, 1)

	d.dialErr = errors.New("dial")
	c.Close()

	nilCount := 0
	errCount := 0
	timeout := time.After(2 * time.Second)
	for i := 0; i < cap(errs); i++ {
		select {
		case err := <-errs:
			switch err {
			case nil:
				nilCount++
			case d.dialErr:
				errCount++
			default:
				t.Fatalf("expected dial error or nil, got %v", err)
			}
		case <-timeout:
			t.Logf("Wait all the time and timeout %d", i)
			return
		}
	}
	if nilCount != 1 {
		t.Errorf("expected one nil error, got %d", nilCount)
	}
	if errCount != cap(errs)-1 {
		t.Errorf("expected %d dial erors, got %d", cap(errs)-1, errCount)
	}
	d.check("done", p, cap(errs), 0)
}
