package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/internal/host"
)

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
	srv.Stop(ctx)
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
	if v := resp.Header.Get("Content-Type"); v != "application/json" {
		t.Fatalf("contentType: %s", v)
	}
	u := new(User)
	if err := json.NewDecoder(resp.Body).Decode(u); err != nil {
		t.Fatal(err)
	}
	if u.Name != "foo" {
		t.Fatalf("got %s want foo", u.Name)
	}
	// POST
	resp, err = http.Post(base+"/users", "application/json", strings.NewReader(`{"name":"bar"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Fatalf("code: %d", resp.StatusCode)
	}
	if v := resp.Header.Get("Content-Type"); v != "application/json" {
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
	req, _ := http.NewRequest("PUT", base+"/users", strings.NewReader(`{"name":"bar"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("code: %d", resp.StatusCode)
	}
	if v := resp.Header.Get("Content-Type"); v != "application/json" {
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
	req, _ = http.NewRequest("OPTIONS", base+"/users", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("code: %d", resp.StatusCode)
	}
	if resp.Header.Get("Access-Control-Allow-Methods") != "OPTIONS" {
		t.Fatal("cors failed")
	}
}
