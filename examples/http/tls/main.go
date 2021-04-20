package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-kratos/kratos/v2"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/home", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "Hello Gorilla Mux!")
	}).Methods("GET")

	httpSrv := transhttp.NewServer(
		transhttp.Address(":8000"),
//		transhttp.CertFile("http/tls/ssl/cert.pem"),
//		transhttp.KeyFile("http/tls/ssl/key.pem"),
	)
	httpSrv.HandlePrefix("/", router)

	app := kratos.New(
		kratos.Name("tls"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Println(err)
	}
}
