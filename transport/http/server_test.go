package http

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/SeeMusic/kratos/v2/errors"
	"github.com/SeeMusic/kratos/v2/middleware"

	"github.com/SeeMusic/kratos/v2/internal/host"
)

type testKey struct{}

type testData struct {
	Path string `json:"path"`
}

func TestServer(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(testData{Path: r.RequestURI})
	}
	ctx := context.Background()
	srv := NewServer()
	srv.HandleFunc("/index", fn)
	srv.HandleFunc("/index/{id:[0-9]+}", fn)
	srv.HandleHeader("content-type", "application/grpc-web+json", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(testData{Path: r.RequestURI})
	})

	if e, err := srv.Endpoint(); err != nil || e == nil || strings.HasSuffix(e.Host, ":0") {
		t.Fatal(e, err)
	}

	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	testHeader(t, srv)
	testClient(t, srv)
	if srv.Stop(ctx) != nil {
		t.Errorf("expected nil got %v", srv.Stop(ctx))
	}
}

func testHeader(t *testing.T, srv *Server) {
	e, err := srv.Endpoint()
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	client, err := NewClient(context.Background(), WithEndpoint(e.Host))
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	reqURL := fmt.Sprintf(e.String() + "/index")
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	req.Header.Set("content-type", "application/grpc-web+json")
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	resp.Body.Close()
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

		{"GET", "/index/1"},
		{"PUT", "/index/1"},
		{"POST", "/index/1"},
		{"PATCH", "/index/1"},
		{"DELETE", "/index/1"},

		{"GET", "/index/notfound"},
	}
	e, err := srv.Endpoint()
	if err != nil {
		t.Fatal(err)
	}
	client, err := NewClient(context.Background(), WithEndpoint(e.Host))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	for _, test := range tests {
		var res testData
		reqURL := fmt.Sprintf(e.String() + test.path)
		req, err := http.NewRequest(test.method, reqURL, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(req)
		if test.path == "/index/notfound" && err != nil {
			if e, ok := err.(*errors.Error); ok && e.Code == http.StatusNotFound {
				continue
			}
		}
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 200 {
			_ = resp.Body.Close()
			t.Fatalf("http status got %d", resp.StatusCode)
		}
		content, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
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
		err := client.Invoke(context.Background(), test.method, test.path, nil, &res)
		if test.path == "/index/notfound" && err != nil {
			if e, ok := err.(*errors.Error); ok && e.Code == http.StatusNotFound {
				continue
			}
		}
		if err != nil {
			t.Fatalf("invoke  error %v", err)
		}
		if res.Path != test.path {
			t.Errorf("expected %s got %s", test.path, res.Path)
		}
	}
}

func BenchmarkServer(b *testing.B) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		data := &testData{Path: r.RequestURI}
		_ = json.NewEncoder(w).Encode(data)
		if r.Context().Value(testKey{}) != "test" {
			w.WriteHeader(500)
		}
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, testKey{}, "test")
	srv := NewServer()
	srv.HandleFunc("/index", fn)
	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	port, ok := host.Port(srv.lis)
	if !ok {
		b.Errorf("expected port got %v", srv.lis)
	}
	client, err := NewClient(context.Background(), WithEndpoint(fmt.Sprintf("127.0.0.1:%d", port)))
	if err != nil {
		b.Errorf("expected nil got %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var res testData
		err := client.Invoke(context.Background(), "POST", "/index", nil, &res)
		if err != nil {
			b.Errorf("expected nil got %v", err)
		}
	}
	_ = srv.Stop(ctx)
}

func TestNetwork(t *testing.T) {
	o := &Server{}
	v := "abc"
	Network(v)(o)
	if !reflect.DeepEqual(v, o.network) {
		t.Errorf("expected %v got %v", v, o.network)
	}
}

func TestAddress(t *testing.T) {
	o := &Server{}
	v := "abc"
	Address(v)(o)
	if !reflect.DeepEqual(v, o.address) {
		t.Errorf("expected %v got %v", v, o.address)
	}
}

func TestTimeout(t *testing.T) {
	o := &Server{}
	v := time.Duration(123)
	Timeout(v)(o)
	if !reflect.DeepEqual(v, o.timeout) {
		t.Errorf("expected %v got %v", v, o.timeout)
	}
}

func TestLogger(t *testing.T) {
	// todo
}

func TestMiddleware(t *testing.T) {
	o := &Server{}
	v := []middleware.Middleware{
		func(middleware.Handler) middleware.Handler { return nil },
	}
	Middleware(v...)(o)
	if !reflect.DeepEqual(v, o.ms) {
		t.Errorf("expected %v got %v", v, o.ms)
	}
}

func TestRequestDecoder(t *testing.T) {
	o := &Server{}
	v := func(*http.Request, interface{}) error { return nil }
	RequestDecoder(v)(o)
	if o.dec == nil {
		t.Errorf("expected nil got %v", o.dec)
	}
}

func TestResponseEncoder(t *testing.T) {
	o := &Server{}
	v := func(http.ResponseWriter, *http.Request, interface{}) error { return nil }
	ResponseEncoder(v)(o)
	if o.enc == nil {
		t.Errorf("expected nil got %v", o.enc)
	}
}

func TestErrorEncoder(t *testing.T) {
	o := &Server{}
	v := func(http.ResponseWriter, *http.Request, error) {}
	ErrorEncoder(v)(o)
	if o.ene == nil {
		t.Errorf("expected nil got %v", o.ene)
	}
}

func TestTLSConfig(t *testing.T) {
	o := &Server{}
	v := &tls.Config{}
	TLSConfig(v)(o)
	if !reflect.DeepEqual(v, o.tlsConf) {
		t.Errorf("expected %v got %v", v, o.tlsConf)
	}
}

func TestListener(t *testing.T) {
	lis := &net.TCPListener{}
	s := &Server{}
	Listener(lis)(s)
	if !reflect.DeepEqual(s.lis, lis) {
		t.Errorf("expected %v got %v", lis, s.lis)
	}
}
