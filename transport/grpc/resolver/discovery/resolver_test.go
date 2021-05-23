package discovery

import (
	"context"
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
}

func (m *testWatch) Next() ([]*registry.ServiceInstance, error) {
	time.Sleep(time.Millisecond * 200)
	ins := []*registry.ServiceInstance{
		{
			ID:      "mock_ID",
			Name:    "mock_Name",
			Version: "mock_Version",
		},
	}
	return ins, nil
}

// Watch creates a watcher according to the service name.
func (m *testWatch) Stop() error {
	return nil
}

func TestWatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	r := &discoveryResolver{
		w:      &testWatch{},
		cc:     &testClientConn{te: t},
		log:    log.NewHelper(log.DefaultLogger),
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
