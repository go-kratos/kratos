package http

import (
	"context"
	"testing"
	"time"
)

type User struct {
	Name string
}

func TestRoute(t *testing.T) {
	ctx := context.Background()
	srv := NewServer()
	route := srv.Route("/")
	route.GET("/users/{name}", func(ctx Context) error {
		u := new(User)
		u.Name = ctx.Vars().Get("name")
		return ctx.Result(200, u)
	})
	route.POST("/users", func(ctx Context) error {
		u := new(User)
		if err := ctx.Bind(u); err != nil {
			return err
		}
		return ctx.Result(200, u)
	})

	go func() {
		if err := srv.Start(ctx); err != nil {
			panic(err)
		}
	}()
	time.Sleep(time.Second)
	srv.Stop(ctx)
}
