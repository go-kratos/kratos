package discovery

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"google.golang.org/grpc/resolver"
)

type testClientConn struct {
	resolver.ClientConn // For unimplemented functions
	te                  *testing.T
}

func (t *testClientConn) UpdateState(s resolver.State) error {
	t.te.Log("UpdateState", s)
	return nil
}

type testWatch struct {
	err error
}

func (m *testWatch) Next() ([]*registry.ServiceInstance, error) {
	time.Sleep(time.Millisecond * 200)
	ins := []*registry.ServiceInstance{
		{
			ID:        "mock_ID",
			Name:      "mock_Name",
			Version:   "mock_Version",
			Endpoints: []string{"grpc://127.0.0.1?isSecure=true"},
		},
		{
			ID:        "mock_ID2",
			Name:      "mock_Name2",
			Version:   "mock_Version2",
			Endpoints: []string{""},
		},
	}
	return ins, m.err
}

// Watch creates a watcher according to the service name.
func (m *testWatch) Stop() error {
	return m.err
}

func TestWatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	r := &discoveryResolver{
		w:        &testWatch{},
		cc:       &testClientConn{te: t},
		log:      log.NewHelper(log.GetLogger()),
		ctx:      ctx,
		cancel:   cancel,
		insecure: false,
	}
	go func() {
		time.Sleep(time.Second * 2)
		r.Close()
	}()
	r.watch()
	t.Log("watch goroutine exited after 2 second")
}

func TestWatchError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	r := &discoveryResolver{
		w:      &testWatch{err: errors.New("bad")},
		cc:     &testClientConn{te: t},
		log:    log.NewHelper(log.GetLogger()),
		ctx:    ctx,
		cancel: cancel,
	}
	go func() {
		time.Sleep(time.Second * 2)
		r.Close()
	}()
	r.watch()
	t.Log("watch goroutine exited after 2 second")
}

func TestWatchContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	r := &discoveryResolver{
		w:      &testWatch{err: context.Canceled},
		cc:     &testClientConn{te: t},
		log:    log.NewHelper(log.GetLogger()),
		ctx:    ctx,
		cancel: cancel,
	}
	go func() {
		time.Sleep(time.Second * 2)
		r.Close()
	}()
	r.watch()
	t.Log("watch goroutine exited after 2 second")
}

func TestParseAttributes(t *testing.T) {
	a := parseAttributes(map[string]string{"a": "b"})
	if !reflect.DeepEqual("b", a.Value("a").(string)) {
		t.Errorf("expect b, got %v", a.Value("a"))
	}
	x := a.WithValue("qq", "ww")
	if !reflect.DeepEqual("ww", x.Value("qq").(string)) {
		t.Errorf("expect ww, got %v", x.Value("qq"))
	}
	if x.Value("notfound") != nil {
		t.Errorf("expect nil, got %v", x.Value("notfound"))
	}
}
