package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/internal/host"
)

const appJSONStr = "application/json"

type User struct {
	Name string `json:"name"`
}

func corsFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			log.Println("cors:", r.Method, r.RequestURI)
			w.Header().Set("Access-Control-Allow-Methods", r.Method)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func authFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("auth:", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func loggingFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("logging:", r.Method, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func TestRoute(t *testing.T) {
	ctx := context.Background()
	srv := NewServer(
		Filter(corsFilter, loggingFilter),
	)
	route := srv.Route("/v1")
	route.GET("/users/{name}", func(ctx Context) error {
		u := new(User)
		u.Name = ctx.Vars().Get("name")
		return ctx.Result(200, u)
	}, authFilter)
	route.POST("/users", func(ctx Context) error {
		u := new(User)
		if err := ctx.Bind(u); err != nil {
			return err
		}
		return ctx.Result(201, u)
	})
	route.PUT("/users", func(ctx Context) error {
		u := new(User)
		if err := ctx.Bind(u); err != nil {
			return err
		}
		h := ctx.Middleware(func(ctx context.Context, in interface{}) (interface{}, error) {
			return u, nil
		})
		return ctx.Returns(h(ctx, u))
	})

	if e, err := srv.Endpoint(); err != nil || e == nil {
		t.Fatal(e, err)
	}
	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	testRoute(t, srv)
	_ = srv.Stop(ctx)
}

func testRoute(t *testing.T, srv *Server) {
	port, ok := host.Port(srv.lis)
	if !ok {
		t.Fatalf("extract port error: %v", srv.lis)
	}
	base := fmt.Sprintf("http://127.0.0.1:%d/v1", port)
	// GET
	resp, err := http.Get(base + "/users/foo")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("code: %d", resp.StatusCode)
	}
	if v := resp.Header.Get("Content-Type"); v != appJSONStr {
		t.Fatalf("contentType: %s", v)
	}
	u := new(User)
	if err = json.NewDecoder(resp.Body).Decode(u); err != nil {
		t.Fatal(err)
	}
	if u.Name != "foo" {
		t.Fatalf("got %s want foo", u.Name)
	}
	// POST
	resp, err = http.Post(base+"/users", appJSONStr, strings.NewReader(`{"name":"bar"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Fatalf("code: %d", resp.StatusCode)
	}
	if v := resp.Header.Get("Content-Type"); v != appJSONStr {
		t.Fatalf("contentType: %s", v)
	}
	u = new(User)
	if err = json.NewDecoder(resp.Body).Decode(u); err != nil {
		t.Fatal(err)
	}
	if u.Name != "bar" {
		t.Fatalf("got %s want bar", u.Name)
	}
	// PUT
	req, _ := http.NewRequest(http.MethodPut, base+"/users", strings.NewReader(`{"name":"bar"}`))
	req.Header.Set("Content-Type", appJSONStr)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("code: %d", resp.StatusCode)
	}
	if v := resp.Header.Get("Content-Type"); v != appJSONStr {
		t.Fatalf("contentType: %s", v)
	}
	u = new(User)
	if err = json.NewDecoder(resp.Body).Decode(u); err != nil {
		t.Fatal(err)
	}
	if u.Name != "bar" {
		t.Fatalf("got %s want bar", u.Name)
	}
	// OPTIONS
	req, _ = http.NewRequest(http.MethodOptions, base+"/users", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("code: %d", resp.StatusCode)
	}
	if resp.Header.Get("Access-Control-Allow-Methods") != http.MethodOptions {
		t.Fatal("cors failed")
	}
}

func TestRouter_Group(t *testing.T) {
	r := &Router{}
	rr := r.Group("a", func(http.Handler) http.Handler { return nil })
	if !reflect.DeepEqual("a", rr.prefix) {
		t.Errorf("expected %q, got %q", "a", rr.prefix)
	}
}

func TestHandle(_ *testing.T) {
	r := newRouter("/", NewServer())
	h := func(i Context) error {
		return nil
	}
	r.GET("/get", h)
	r.HEAD("/head", h)
	r.PATCH("/patch", h)
	r.DELETE("/delete", h)
	r.CONNECT("/connect", h)
	r.OPTIONS("/options", h)
	r.TRACE("/trace", h)
}

func TestRouter_ContextDataRace(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx := context.Background()
	signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL)

	srv := NewServer(Address(":18080"), Timeout(time.Millisecond*50))
	router := srv.Route("/")
	router.GET("/ping", func(ctx Context) error {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://www.baidu.com", nil)
		_, _ = http.DefaultClient.Do(req)
		return ctx.String(200, "pong")
	})

	// start server
	go func() {
		if err := srv.Start(ctx); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			t.Error(err)
			return
		}
	}()

	time.Sleep(time.Second)

	// start client
	workers := 100
	for i := 0; i < workers; i++ {
		_i := i
		go func() {
			for {
				req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:18080/ping", nil)
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					break
				}
				_ = res.Body.Close()
			}
			t.Logf("worker: %d shutdown\n", _i)
		}()
	}

	<-ctx.Done()
	t.Log("test end")
}
