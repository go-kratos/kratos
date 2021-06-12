package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func main() {
	router := gin.Default()
	router.GET("/home", func(ctx *gin.Context) {
		ctx.String(200, "Hello Gin!")
	})

	httpSrv := http.NewServer(http.Address(":8000"))
	httpSrv.HandlePrefix("/", router)

	app := kratos.New(
		kratos.Name("gin"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
