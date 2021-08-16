package main

import (
	"log"

	"github.com/go-kratos/kratos/examples/http/query_array/hello"
	"github.com/go-kratos/kratos/examples/http/query_array/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func main() {
	s := service.NewGreeterService()
	httpSrv := http.NewServer(
		http.Address(":8080"),
		http.Middleware(
			recovery.Recovery(),
		),
	)
	hello.RegisterGreeterHTTPServer(httpSrv, s)

	app := kratos.New(
		kratos.Server(
			httpSrv,
		),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
