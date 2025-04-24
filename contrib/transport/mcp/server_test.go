package mcp

import (
	"context"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	var (
		ctx = context.Background()
		srv = NewServer("test", "v1.0.0", Address(":0"))
	)
	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	if err := srv.Stop(ctx); err != nil {
		t.Fatal(err)
	}
}
