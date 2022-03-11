package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SeeMusic/kratos/v2"
	transhttp "github.com/SeeMusic/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/home", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "Hello Gorilla Mux!")
	}).Methods("GET")

	httpSrv := transhttp.NewServer(transhttp.Address(":8000"))
	httpSrv.HandlePrefix("/", router)

	app := kratos.New(
		kratos.Name("mux"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
