package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	if e, err := srv.Endpoint(); err != nil || e == "" {
		t.Fatal(e, err)
	}

	go func() {
		if err := srv.Start(); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	testClient(t, srv)
	srv.Stop()
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
	port, ok := host.Port(srv.lis)
	if !ok {
		t.Fatalf("extract port error: %v", srv.lis)
	}
	client, err := NewClient(context.Background(), WithEndpoint(fmt.Sprintf("127.0.0.1:%d", port)))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range tests {
		var res testData
		url := fmt.Sprintf("http://127.0.0.1:%d%s", port, test.path)
		req, err := http.NewRequest(test.method, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read resp error %v", err)
		}
		err = json.Unmarshal(content, &res)
		if err != nil {
			t.Fatalf("unmarshal resp error %v", err)
		}
		if res.Path != test.path {
			t.Errorf("expected %s got %s", test.path, res.Path)
		}
	}
	for _, test := range tests {
		var res testData
		err := client.Invoke(context.Background(), test.path, nil, &res, Method(test.method))
		if err != nil {
			t.Fatalf("invoke  error %v", err)
		}
		if res.Path != test.path {
			t.Errorf("expected %s got %s", test.path, res.Path)
		}
	}

}
