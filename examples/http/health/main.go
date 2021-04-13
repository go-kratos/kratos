package main

import (
	"context"
	"errors"
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/health"
)

type mysqlChecker struct{}

func (c *mysqlChecker) CheckHealth(ctx context.Context) error {
	return nil
}

type redisChecker struct{}

func (c *redisChecker) CheckHealth(ctx context.Context) error {
	return errors.New("connection refused")
}

func main() {
	handler := health.NewHandler()
	handler.AddChecker("mysql", &mysqlChecker{})
	handler.AddObserver("redis", &redisChecker{})

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
