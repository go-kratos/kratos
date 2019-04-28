package breaker

import (
	"errors"
	"testing"
	"time"

	xtime "go-common/library/time"
)

func TestGroup(t *testing.T) {
	g1 := NewGroup(nil)
	g2 := NewGroup(_conf)
	if g1.conf != g2.conf {
		t.FailNow()
	}

	brk := g2.Get("key")
	brk1 := g2.Get("key1")
	if brk == brk1 {
		t.FailNow()
	}
	brk2 := g2.Get("key")
	if brk != brk2 {
		t.FailNow()
	}

	g := NewGroup(_conf)
	c := &Config{
		Window:    xtime.Duration(1 * time.Second),
		Sleep:     xtime.Duration(100 * time.Millisecond),
		Bucket:    10,
		Ratio:     0.5,
		Request:   100,
		SwitchOff: !_conf.SwitchOff,
	}
	g.Reload(c)
	if g.conf.SwitchOff == _conf.SwitchOff {
		t.FailNow()
	}
}

func TestInit(t *testing.T) {
	switchOff := _conf.SwitchOff
	c := &Config{
		Window:    xtime.Duration(3 * time.Second),
		Sleep:     xtime.Duration(100 * time.Millisecond),
		Bucket:    10,
		Ratio:     0.5,
		Request:   100,
		SwitchOff: !switchOff,
	}
	Init(c)
	if _conf.SwitchOff == switchOff {
		t.FailNow()
	}
}

func TestGo(t *testing.T) {
	if err := Go("test_run", func() error {
		t.Log("breaker allow,callback run()")
		return nil
	}, func() error {
		t.Log("breaker not allow,callback fallback()")
		return errors.New("breaker not allow")
	}); err != nil {
		t.Error(err)
	}

	_group.Reload(&Config{
		Window:    xtime.Duration(3 * time.Second),
		Sleep:     xtime.Duration(100 * time.Millisecond),
		Bucket:    10,
		Ratio:     0.5,
		Request:   100,
		SwitchOff: true,
	})

	if err := Go("test_fallback", func() error {
		t.Log("breaker allow,callback run()")
		return nil
	}, func() error {
		t.Log("breaker not allow,callback fallback()")
		return nil
	}); err != nil {
		t.Error(err)
	}
}

func markSuccess(b Breaker, count int) {
	for i := 0; i < count; i++ {
		b.MarkSuccess()
	}
}

func markFailed(b Breaker, count int) {
	for i := 0; i < count; i++ {
		b.MarkFailed()
	}
}
