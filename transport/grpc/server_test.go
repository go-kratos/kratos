package grpc

import (
	"context"
	"crypto/tls"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"google.golang.org/grpc"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testKey struct{}

func TestServer(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, testKey{}, "test")
	srv := NewServer(Middleware([]middleware.Middleware{
		func(middleware.Handler) middleware.Handler { return nil },
	}...))

	if e, err := srv.Endpoint(); err != nil || e == nil || strings.HasSuffix(e.Host, ":0") {
		t.Fatal(e, err)
	}

	go func() {
		// start server
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	testClient(t, srv)
	srv.Stop(ctx)
}

func testClient(t *testing.T, srv *Server) {
	u, err := srv.Endpoint()
	if err != nil {
		t.Fatal(err)
	}
	// new a gRPC client
	conn, err := DialInsecure(context.Background(), WithEndpoint(u.Host))
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
}

func TestNetwork(t *testing.T) {
	o := &Server{}
	v := "abc"
	Network(v)(o)
	assert.Equal(t, v, o.network)
}

func TestAddress(t *testing.T) {
	o := &Server{}
	v := "abc"
	Address(v)(o)
	assert.Equal(t, v, o.address)
}

func TestTimeout(t *testing.T) {
	o := &Server{}
	v := time.Duration(123)
	Timeout(v)(o)
	assert.Equal(t, v, o.timeout)
}

func TestMiddleware(t *testing.T) {
	o := &Server{}
	v := []middleware.Middleware{
		func(middleware.Handler) middleware.Handler { return nil },
	}
	Middleware(v...)(o)
	assert.Equal(t, v, o.middleware)
}

type mockLogger struct {
	level log.Level
	key   string
	val   string
}

func (l *mockLogger) Log(level log.Level, keyvals ...interface{}) error {
	l.level = level
	l.key = keyvals[0].(string)
	l.val = keyvals[1].(string)
	return nil
}

func TestLogger(t *testing.T) {
	o := &Server{}
	v := &mockLogger{}
	Logger(v)(o)
	o.log.Log(log.LevelWarn, "foo", "bar")
	assert.Equal(t, "foo", v.key)
	assert.Equal(t, "bar", v.val)
	assert.Equal(t, log.LevelWarn, v.level)
}

func TestTLSConfig(t *testing.T) {
	o := &Server{}
	v := &tls.Config{}
	TLSConfig(v)(o)
	assert.Equal(t, v, o.tlsConf)
}

func TestUnaryInterceptor(t *testing.T) {
	o := &Server{}
	v := []grpc.UnaryServerInterceptor{
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return nil, nil
		},
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return nil, nil
		},
	}
	UnaryInterceptor(v...)(o)
	assert.Equal(t, v, o.ints)
	}

func TestOptions(t *testing.T) {
	o := &Server{}
	v := []grpc.ServerOption{
		grpc.EmptyServerOption{},
	}
	Options(v...)(o)
	assert.Equal(t, v, o.grpcOpts)
}
