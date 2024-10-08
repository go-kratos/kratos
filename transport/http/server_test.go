package http

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	kratoserrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/internal/host"
)

var h = func(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(testData{Path: r.RequestURI})
}

type testKey struct{}

type testData struct {
	Path string `json:"path"`
}

// handleFuncWrapper is a wrapper for http.HandlerFunc to implement http.Handler
type handleFuncWrapper struct {
	fn http.HandlerFunc
}

func (x *handleFuncWrapper) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	x.fn.ServeHTTP(writer, request)
}

func newHandleFuncWrapper(fn http.HandlerFunc) http.Handler {
	return &handleFuncWrapper{fn: fn}
}

func TestServeHTTP(t *testing.T) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	mux := NewServer(Listener(ln))
	mux.HandleFunc("/index", h)
	mux.Route("/errors").GET("/cause", func(Context) error {
		return kratoserrors.BadRequest("xxx", "zzz").
			WithMetadata(map[string]string{"foo": "bar"}).
			WithCause(errors.New("error cause"))
	})
	if err = mux.WalkRoute(func(r RouteInfo) error {
		t.Logf("WalkRoute: %+v", r)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if e, err := mux.Endpoint(); err != nil || e == nil || strings.HasSuffix(e.Host, ":0") {
		t.Fatal(e, err)
	}
	srv := http.Server{Handler: mux}
	go func() {
		if err := srv.Serve(ln); err != nil {
			if kratoserrors.Is(err, http.ErrServerClosed) {
				return
			}
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	if err := srv.Shutdown(context.Background()); err != nil {
		t.Log(err)
	}
}

func TestServer(t *testing.T) {
	ctx := context.Background()
	srv := NewServer()
	srv.Handle("/index", newHandleFuncWrapper(h))
	srv.HandleFunc("/index/{id:[0-9]+}", h)
	srv.HandlePrefix("/test/prefix", newHandleFuncWrapper(h))
	srv.HandleHeader("content-type", "application/grpc-web+json", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(testData{Path: r.RequestURI})
	})
	srv.Route("/errors").GET("/cause", func(Context) error {
		return kratoserrors.BadRequest("xxx", "zzz").
			WithMetadata(map[string]string{"foo": "bar"}).
			WithCause(errors.New("error cause"))
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
	testAccept(t, srv)
	time.Sleep(time.Second)
	if srv.Stop(ctx) != nil {
		t.Errorf("expected nil got %v", srv.Stop(ctx))
	}
}

func testAccept(t *testing.T, srv *Server) {
	tests := []struct {
		method      string
		path        string
		contentType string
	}{
		{http.MethodGet, "/errors/cause", "application/json"},
		{http.MethodGet, "/errors/cause", "application/proto"},
	}
	e, err := srv.Endpoint()
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	client, err := NewClient(context.Background(), WithEndpoint(e.Host))
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	for _, test := range tests {
		req, err := http.NewRequest(test.method, e.String()+test.path, nil)
		if err != nil {
			t.Errorf("expected nil got %v", err)
		}
		req.Header.Set("Content-Type", test.contentType)
		resp, err := client.Do(req)
		if kratoserrors.Code(err) != 400 {
			t.Errorf("expected 400 got %v", err)
		}
		if err == nil {
			resp.Body.Close()
		}
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
	req, err := http.NewRequest(http.MethodGet, e.String()+"/index", nil)
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
		code   int
	}{
		{http.MethodGet, "/index", http.StatusOK},
		{http.MethodPut, "/index", http.StatusOK},
		{http.MethodPost, "/index", http.StatusOK},
		{http.MethodPatch, "/index", http.StatusOK},
		{http.MethodDelete, "/index", http.StatusOK},

		{http.MethodGet, "/index/1", http.StatusOK},
		{http.MethodPut, "/index/1", http.StatusOK},
		{http.MethodPost, "/index/1", http.StatusOK},
		{http.MethodPatch, "/index/1", http.StatusOK},
		{http.MethodDelete, "/index/1", http.StatusOK},

		{http.MethodGet, "/index/notfound", http.StatusNotFound},
		{http.MethodGet, "/errors/cause", http.StatusBadRequest},
		{http.MethodGet, "/test/prefix/123111", http.StatusOK},
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
		reqURL := e.String() + test.path
		req, err := http.NewRequest(test.method, reqURL, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(req)
		if kratoserrors.Code(err) != test.code {
			t.Fatalf("want %v, but got %v", test, err)
		}
		if err != nil {
			continue
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
		if kratoserrors.Code(err) != test.code {
			t.Fatalf("want %v, but got %v", test, err)
		}
		if err != nil {
			continue
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
		err := client.Invoke(context.Background(), http.MethodPost, "/index", nil, &res)
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

func TestRequestDecoder(t *testing.T) {
	o := &Server{}
	v := func(*http.Request, interface{}) error { return nil }
	RequestDecoder(v)(o)
	if o.decBody == nil {
		t.Errorf("expected nil got %v", o.decBody)
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

func TestStrictSlash(t *testing.T) {
	o := &Server{}
	v := true
	StrictSlash(v)(o)
	if !reflect.DeepEqual(v, o.strictSlash) {
		t.Errorf("expected %v got %v", v, o.tlsConf)
	}
}

func TestListener(t *testing.T) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	s := &Server{}
	Listener(lis)(s)
	if !reflect.DeepEqual(s.lis, lis) {
		t.Errorf("expected %v got %v", lis, s.lis)
	}
	if e, err := s.Endpoint(); err != nil || e == nil {
		t.Errorf("expected not empty")
	}
}

func TestNotFoundHandler(t *testing.T) {
	mux := http.NewServeMux()
	srv := NewServer(NotFoundHandler(mux))
	if !reflect.DeepEqual(srv.router.NotFoundHandler, mux) {
		t.Errorf("expected %v got %v", mux, srv.router.NotFoundHandler)
	}
}

func TestMethodNotAllowedHandler(t *testing.T) {
	mux := http.NewServeMux()
	srv := NewServer(MethodNotAllowedHandler(mux))
	if !reflect.DeepEqual(srv.router.MethodNotAllowedHandler, mux) {
		t.Errorf("expected %v got %v", mux, srv.router.MethodNotAllowedHandler)
	}
}
