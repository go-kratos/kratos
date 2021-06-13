package main

import (
	"log"

	"github.com/go-kratos/kratos/examples/ws/handler"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ws", handler.WsHandler)

	httpSrv := http.NewServer(http.Address(":8080"))
	httpSrv.HandlePrefix("/", router)

	app := kratos.New(
		kratos.Name("ws"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Println(err)
	}
}
