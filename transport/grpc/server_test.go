package grpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/internal/host"
)

func TestServer(t *testing.T) {
	srv := NewServer()
	if endpoint, err := srv.Endpoint(); err != nil || endpoint == "" {
		t.Fatal(endpoint, err)
	}

	time.AfterFunc(time.Second, func() {
		defer srv.Stop()
		testClient(t, srv)
	})
	// start server
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
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
