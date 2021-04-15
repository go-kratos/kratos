package main

import (
	"context"
	"errors"
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/health"
)

func main() {
	handler := health.NewHandler()
	handler.AddChecker("mysql", func(ctx context.Context) error {
		return nil
	})
	handler.AddObserver("redis", func(ctx context.Context) error {
		return errors.New("connection refused")
	})

	httpSrv := http.NewServer(http.Address(":8000"))
	httpSrv.Handle("/healthz", handler)

	app := kratos.New(
		kratos.Name("mux"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Println(err)
	}
}
