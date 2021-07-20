package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
)

func hello(ctx http.Context) error {
	name := ctx.Vars().Get("name")
	return ctx.String(200, "hellowolrd "+name)
}

func main() {
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Filter(handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST"}),
		)),
	)
	route := httpSrv.Route("/")
	route.GET("/helloworld/{name}", hello)

	app := kratos.New(
		kratos.Name("cors"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
