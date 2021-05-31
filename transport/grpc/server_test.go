package grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/internal/host"
)

type testKey struct{}

func TestServer(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, testKey{}, "test")
	srv := NewServer()
	if e, err := srv.Endpoint(); err != nil || e == nil {
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
	port, ok := host.Port(srv.lis)
	if !ok {
		t.Fatalf("extract port error: %v", srv.lis)
	}
	endpoint := fmt.Sprintf("127.0.0.1:%d", port)
	// new a gRPC client
	conn, err := DialInsecure(context.Background(), WithEndpoint(endpoint))
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
}
