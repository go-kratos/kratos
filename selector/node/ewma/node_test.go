package ewma

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/registry"
	"github.com/go-kratos/kratos/v3/selector"
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

	time.Sleep(time.Millisecond * 15)
	done(context.Background(), selector.DoneInfo{})
	if float64(70) >= wn.Weight() {
		t.Errorf("float64(30000) >= wn.Weight()(%v)", wn.Weight())
	}
	if float64(1200) <= wn.Weight() {
		t.Errorf("float64(1000) <= wn.Weight()(%v)", wn.Weight())
	}
	if time.Millisecond*30 <= wn.PickElapsed() {
		t.Errorf("time.Millisecond*30 <= wn.PickElapsed()(%v)", wn.PickElapsed())
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
	if float64(1000) >= wn.Weight() {
		t.Errorf("float64(1000) >= wn.Weight()(%v)", wn.Weight())
	}
	if float64(2000) <= wn.Weight() {
		t.Errorf("float64(2000) <= wn.Weight()(%v)", wn.Weight())
	}
}

// TestCanceledDoesNotDegradeNode checks the health signal for each error
// kind. Cancellation, direct or wrapped by a transport, must not lower node
// health; backend failures still must. Each case uses a fresh node and a
// single completion: with no previous lag sample the EWMA decay factor is
// zero, so success is exactly 1000 or 0 and the test needs no sleeps.
func TestCanceledDoesNotDegradeNode(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		errHandler func(error) bool
		want       uint64
	}{
		{name: "nil", err: nil, want: 1000},
		{name: "canceled", err: context.Canceled, want: 1000},
		{name: "wrapped canceled", err: fmt.Errorf("call: %w", context.Canceled), want: 1000},
		{name: "http canceled", err: &url.Error{Op: "Get", URL: "http://127.0.0.1:9090", Err: context.Canceled}, want: 1000},
		{name: "dial canceled", err: &net.OpError{Op: "dial", Net: "tcp", Err: context.Canceled}, want: 1000},
		{name: "deadline exceeded", err: context.DeadlineExceeded, want: 0},
		{name: "wrapped deadline exceeded", err: fmt.Errorf("call: %w", context.DeadlineExceeded), want: 0},
		{name: "http deadline exceeded", err: &url.Error{Op: "Get", URL: "http://127.0.0.1:9090", Err: context.DeadlineExceeded}, want: 0},
		{name: "net error", err: &net.OpError{Op: "dial", Net: "tcp", Err: fmt.Errorf("connection refused")}, want: 0},
		{name: "service unavailable", err: errors.ServiceUnavailable("", ""), want: 0},
		{name: "gateway timeout", err: errors.GatewayTimeout("", ""), want: 0},
		{name: "service unavailable canceled while reading", err: errors.ServiceUnavailable("", "").WithCause(context.Canceled), want: 0},
		{name: "gateway timeout canceled while reading", err: errors.GatewayTimeout("", "").WithCause(context.Canceled), want: 0},
		{name: "other error", err: fmt.Errorf("bad request"), want: 1000},
		{
			name:       "error handler can still penalize canceled",
			err:        context.Canceled,
			errHandler: func(err error) bool { return errors.Is(err, context.Canceled) },
			want:       0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := (&Builder{ErrHandler: tt.errHandler}).Build(nil).(*Node)
			n.Pick()(context.Background(), selector.DoneInfo{Err: tt.err})
			if got := n.success.Load(); got != tt.want {
				t.Errorf("success = %d, want %d", got, tt.want)
			}
		})
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
	errs := []error{
		context.DeadlineExceeded,
		context.Canceled,
		net.ErrClosed,
	}
	for i := 0; i < 5; i++ {
		var err error
		if i != 0 {
			err = errs[i%len(errs)]
		}
		done := wn.Pick()
		if done == nil {
			t.Errorf("expect not nil, got nil")
		}
		time.Sleep(time.Millisecond * 20)
		done(context.Background(), selector.DoneInfo{Err: err})
	}
	if float64(1000) >= wn.Weight() {
		t.Errorf("float64(100) >= wn.Weight()(%v)", wn.Weight())
	}
	if float64(2000) <= wn.Weight() {
		t.Errorf("float64(200) <= wn.Weight()(%v)", wn.Weight())
	}
}

func BenchmarkPickAndWeight(b *testing.B) {
	bu := &Builder{}
	node := bu.Build(selector.NewNode(
		"http",
		"127.0.0.1:9090",
		&registry.ServiceInstance{
			ID:        "127.0.0.1:9090",
			Name:      "helloworld",
			Version:   "v1.0.0",
			Endpoints: []string{"http://127.0.0.1:9090"},
			Metadata:  map[string]string{"weight": "10"},
		}))
	di := selector.DoneInfo{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			done := node.Pick()
			node.Weight()
			done(context.Background(), di)
		}
	})
}
