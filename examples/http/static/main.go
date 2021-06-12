package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/go-kratos/kratos/v2"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

//go:embed assets/*
var f embed.FS

func main() {
	router := mux.NewRouter()
	// example: /assets/index.html
	router.PathPrefix("/assets").Handler(http.FileServer(http.FS(f)))

	httpSrv := transhttp.NewServer(transhttp.Address(":8000"))
	httpSrv.HandlePrefix("/", router)

	app := kratos.New(
		kratos.Name("static"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
