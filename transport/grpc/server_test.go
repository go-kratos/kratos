package grpc

import (
	"context"
	"crypto/tls"
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
	srv := NewServer()

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
	o := &clientOptions{}
	v := []middleware.Middleware{
		func(middleware.Handler) middleware.Handler { return nil },
	}
	WithMiddleware(v...)(o)
	assert.Equal(t, v, o.middleware)
}

func TestLogger(t *testing.T) {
	//todo
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
