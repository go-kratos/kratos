package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {
	router := echo.New()
	router.GET("/home", func(context echo.Context) error {
		return context.JSON(200,"Hello echo")
	})

	httpSrv := http.NewServer(http.Address(":9527"))
	httpSrv.HandlePrefix("/", router)

	app := kratos.New(
		kratos.Name("echo"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}