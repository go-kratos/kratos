package kratos

import (
	"context"
	"log"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	xlog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport"
)

func TestID(t *testing.T) {
	o := &options{}
	v := "123"
	ID(v)(o)
	if !reflect.DeepEqual(v, o.id) {
		t.Fatalf("o.id:%s is not equal to v:%s", o.id, v)
	}
}

func TestName(t *testing.T) {
	o := &options{}
	v := "abc"
	Name(v)(o)
	if !reflect.DeepEqual(v, o.name) {
		t.Fatalf("o.name:%s is not equal to v:%s", o.name, v)
	}
}

func TestVersion(t *testing.T) {
	o := &options{}
	v := "123"
	Version(v)(o)
	if !reflect.DeepEqual(v, o.version) {
		t.Fatalf("o.version:%s is not equal to v:%s", o.version, v)
	}
}

func TestMetadata(t *testing.T) {
	o := &options{}
	v := map[string]string{
		"a": "1",
		"b": "2",
	}
	Metadata(v)(o)
	if !reflect.DeepEqual(v, o.metadata) {
		t.Fatalf("o.metadata:%s is not equal to v:%s", o.metadata, v)
	}
}

func TestEndpoint(t *testing.T) {
	o := &options{}
	v := []*url.URL{
		{Host: "example.com"},
		{Host: "foo.com"},
	}
	Endpoint(v...)(o)
	if !reflect.DeepEqual(v, o.endpoints) {
		t.Fatalf("o.endpoints:%s is not equal to v:%s", o.endpoints, v)
	}
}

func TestContext(t *testing.T) {
	type ctxKey struct {
		Key string
	}
	o := &options{}
	v := context.WithValue(context.TODO(), ctxKey{Key: "context"}, "b")
	Context(v)(o)
	if !reflect.DeepEqual(v, o.ctx) {
		t.Fatalf("o.ctx:%s is not equal to v:%s", o.ctx, v)
	}
}

func TestLogger(t *testing.T) {
	o := &options{}
	v := xlog.NewStdLogger(log.Writer())
	Logger(v)(o)
	if !reflect.DeepEqual(v, o.logger) {
		t.Fatalf("o.logger:%v is not equal to xlog.NewHelper(v):%v", o.logger, xlog.NewHelper(v))
	}
}

type mockServer struct{}

func (m *mockServer) Start(_ context.Context) error { return nil }
func (m *mockServer) Stop(_ context.Context) error  { return nil }

func TestServer(t *testing.T) {
	o := &options{}
	v := []transport.Server{
		&mockServer{}, &mockServer{},
	}
	Server(v...)(o)
	if !reflect.DeepEqual(v, o.servers) {
		t.Fatalf("o.servers:%s is not equal to xlog.NewHelper(v):%s", o.servers, v)
	}
}

type mockSignal struct{}

func (m *mockSignal) String() string { return "sig" }
func (m *mockSignal) Signal()        {}

func TestSignal(t *testing.T) {
	o := &options{}
	v := []os.Signal{
		&mockSignal{}, &mockSignal{},
	}
	Signal(v...)(o)
	if !reflect.DeepEqual(v, o.sigs) {
		t.Fatal("o.sigs is not equal to v")
	}
}

type mockRegistrar struct{}

func (m *mockRegistrar) Register(_ context.Context, _ *registry.ServiceInstance) error {
	return nil
}

func (m *mockRegistrar) Deregister(_ context.Context, _ *registry.ServiceInstance) error {
	return nil
}

func TestRegistrar(t *testing.T) {
	o := &options{}
	v := &mockRegistrar{}
	Registrar(v)(o)
	if !reflect.DeepEqual(v, o.registrar) {
		t.Fatal("o.registrar is not equal to v")
	}
}

func TestRegistrarTimeout(t *testing.T) {
	o := &options{}
	v := time.Duration(123)
	RegistrarTimeout(v)(o)
	if !reflect.DeepEqual(v, o.registrarTimeout) {
		t.Fatal("o.registrarTimeout is not equal to v")
	}
}

func TestStopTimeout(t *testing.T) {
	o := &options{}
	v := time.Duration(123)
	StopTimeout(v)(o)
	if !reflect.DeepEqual(v, o.stopTimeout) {
		t.Fatal("o.stopTimeout is not equal to v")
	}
}

func TestBeforeStart(t *testing.T) {
	o := &options{}
	v := func(_ context.Context) error {
		t.Log("BeforeStart...")
		return nil
	}
	BeforeStart(v)(o)
}

func TestBeforeStop(t *testing.T) {
	o := &options{}
	v := func(_ context.Context) error {
		t.Log("BeforeStop...")
		return nil
	}
	BeforeStop(v)(o)
}

func TestAfterStart(t *testing.T) {
	o := &options{}
	v := func(_ context.Context) error {
		t.Log("AfterStart...")
		return nil
	}
	AfterStart(v)(o)
}

func TestAfterStop(t *testing.T) {
	o := &options{}
	v := func(_ context.Context) error {
		t.Log("AfterStop...")
		return nil
	}
	AfterStop(v)(o)
}
