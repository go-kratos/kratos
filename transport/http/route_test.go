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

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("auth:", r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	}
}
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("logging:", r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	}
}

func TestRoute(t *testing.T) {
	ctx := context.Background()
	srv := NewServer(
		RouteMiddleware(loggingMiddleware),
	)
	route := srv.Route("/")
	route.GET("/users/{name}", func(ctx Context) error {
		u := new(User)
		u.Name = ctx.Vars().Get("name")
		return ctx.Result(200, u)
	}, authMiddleware)
	route.POST("/users", func(ctx Context) error {
		u := new(User)
		if err := ctx.Bind(u); err != nil {
			return err
		}
		return ctx.Result(201, u)
	})

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
	base := fmt.Sprintf("http://127.0.0.1:%d", port)
	// GET
	resp, err := http.Get(base + "/users/foo")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("code: %d", resp.StatusCode)
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
	u = new(User)
	if err = json.NewDecoder(resp.Body).Decode(u); err != nil {
		t.Fatal(err)
	}
	if u.Name != "bar" {
		t.Fatalf("got %s want bar", u.Name)
	}
}
