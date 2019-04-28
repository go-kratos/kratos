package supervisor

import (
	"testing"
	"time"
)

func create() *Supervisor {
	now := time.Now()
	end := now.Add(time.Hour * 1)
	conf := &Config{
		On:    true,
		Begin: now,
		End:   end,
	}
	return New(conf)
}

func TestSupervisor(t *testing.T) {
	sv := create()
	in := sv.conf.Begin.Add(time.Second * 10)
	out := sv.conf.End.Add(time.Second * 10)

	if sv.forbid("GET", in) {
		t.Error("Request should never be blocked on GET method")
	}

	if !sv.forbid("POST", in) {
		t.Errorf("Request should be blocked on POST method at %+v", in)
	}

	if sv.forbid("POST", out) {
		t.Errorf("Request should not be blocked at %+v", out)
	}
}

func TestReload(t *testing.T) {
	zero := time.Unix(0, 0)
	conf := &Config{
		On:    false,
		Begin: zero,
		End:   zero,
	}
	sv := create()

	// reload with nil
	sv.Reload(nil)

	// reload with valid config
	sv.Reload(conf)

	if sv.conf != conf && sv.on == false {
		t.Errorf("Failed to reload config %+v, current config is %+v", conf, sv.conf)
	}
}
