package ewma

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/selector"
)

func TestDirect(t *testing.T) {
	b := &Builder{}
	wn := b.Build(selector.NewNode(
		"http",
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
			Metadata:  map[string]string{"weight": "10"},
		}))

	if !reflect.DeepEqual(float64(100), wn.Weight()) {
		t.Errorf("expect %v, got %v", 100, wn.Weight())
	}
	done := wn.Pick()
	if done == nil {
		t.Errorf("done is equal to nil")
	}
	done2 := wn.Pick()
	if done2 == nil {
		t.Errorf("done2 is equal to nil")
	}

	time.Sleep(time.Millisecond * 10)
	done(context.Background(), selector.DoneInfo{})
	if float64(30000) >= wn.Weight() {
		t.Errorf("float64(30000) >= wn.Weight()(%v)", wn.Weight())
	}
	if float64(60000) <= wn.Weight() {
		t.Errorf("float64(60000) <= wn.Weight()(%v)", wn.Weight())
	}
	if time.Millisecond*15 <= wn.PickElapsed() {
		t.Errorf("time.Millisecond*15 <= wn.PickElapsed()(%v)", wn.PickElapsed())
	}
	if time.Millisecond*5 >= wn.PickElapsed() {
		t.Errorf("time.Millisecond*5 >= wn.PickElapsed()(%v)", wn.PickElapsed())
	}
}

func TestDirectError(t *testing.T) {
	b := &Builder{}
	wn := b.Build(selector.NewNode(
		"http",
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
			Metadata:  map[string]string{"weight": "10"},
		}))

	for i := 0; i < 5; i++ {
		var err error
		if i != 0 {
			err = context.DeadlineExceeded
		}
		done := wn.Pick()
		if done == nil {
			t.Errorf("expect not nil, got nil")
		}
		time.Sleep(time.Millisecond * 20)
		done(context.Background(), selector.DoneInfo{Err: err})
	}
	if float64(30000) >= wn.Weight() {
		t.Errorf("float64(30000) >= wn.Weight()(%v)", wn.Weight())
	}
	if float64(60000) <= wn.Weight() {
		t.Errorf("float64(60000) <= wn.Weight()(%v)", wn.Weight())
	}
}

func TestDirectErrorHandler(t *testing.T) {
	b := &Builder{
		ErrHandler: func(err error) bool {
			return err != nil
		},
	}
	wn := b.Build(selector.NewNode(
		"http",
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
			Metadata:  map[string]string{"weight": "10"},
		}))

	for i := 0; i < 5; i++ {
		var err error
		if i != 0 {
			err = context.DeadlineExceeded
		}
		done := wn.Pick()
		if done == nil {
			t.Errorf("expect not nil, got nil")
		}
		time.Sleep(time.Millisecond * 20)
		done(context.Background(), selector.DoneInfo{Err: err})
	}
	if float64(30000) >= wn.Weight() {
		t.Errorf("float64(30000) >= wn.Weight()(%v)", wn.Weight())
	}
	if float64(60000) <= wn.Weight() {
		t.Errorf("float64(60000) <= wn.Weight()(%v)", wn.Weight())
	}
}
