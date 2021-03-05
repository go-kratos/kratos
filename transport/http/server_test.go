package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/internal/host"
)

type testData struct {
	Path string `json:"path"`
}

func TestServer(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		data := &testData{Path: r.RequestURI}
		json.NewEncoder(w).Encode(data)
	}
	srv := NewServer()
	srv.HandleFunc("/index", fn)

	time.AfterFunc(time.Second, func() {
		defer srv.Stop()
		testClient(t, srv)
	})

	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}
}

func testClient(t *testing.T, srv *Server) {
	tests := []struct {
		method string
		path   string
	}{
		{"GET", "/index"},
		{"PUT", "/index"},
		{"POST", "/index"},
		{"PATCH", "/index"},
		{"DELETE", "/index"},
	}
	client, err := NewClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	port, ok := host.Port(srv.lis)
	if !ok {
		t.Fatalf("extract port error: %v", srv.lis)
	}
	for _, test := range tests {
		var res testData
		url := fmt.Sprintf("http://127.0.0.1:%d%s", port, test.path)
		req, err := http.NewRequest(test.method, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := Do(client, req, &res); err != nil {
			t.Fatal(err)
		}
		if res.Path != test.path {
			t.Errorf("expected %s got %s", test.path, res.Path)
		}
	}
}
