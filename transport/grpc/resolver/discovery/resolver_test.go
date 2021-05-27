package discovery

import (
	"context"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"google.golang.org/grpc/resolver"
)

type service struct {
	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []string
}

func (s *service) ID() string {
	return s.id
}

func (s *service) Name() string {
	return s.name
}

func (s *service) Version() string {
	return s.version
}

func (s *service) Metadata() map[string]string {
	return s.metadata
}

func (s *service) Endpoints() []string {
	return s.endpoints
}

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

func (m *testWatch) Next() ([]registry.Service, error) {
	time.Sleep(time.Millisecond * 200)
	var inss []registry.Service = []registry.Service{&service{
		id:      "mock_ID",
		name:    "mock_Name",
		version: "mock_Version",
	}}

	return inss, nil
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
