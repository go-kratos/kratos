package main

import (
	_ "embed"
	"fmt"
	"github.com/go-kratos/kratos/v2"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

//go:embed cert.pem
var cert []byte

//go:embed key.pem
var key []byte

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/home", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "Hello Gorilla Mux!")
	}).Methods("GET")

	httpSrv := transhttp.NewServer(
		transhttp.Address(":8000"),
		transhttp.X509KeyPair(cert,key),
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
