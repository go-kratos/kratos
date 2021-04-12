package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-kratos/kratos/v2"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.FormValue("name")
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, _ = io.Copy(f, file)
	fmt.Fprint(w, "File "+fileName+" Uploaded successfully")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/upload", uploadFile).Methods("POST")

	httpSrv := transhttp.NewServer(transhttp.Address(":8000"))
	httpSrv.HandlePrefix("/", router)

	app := kratos.New(
		kratos.Name("upload"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Println(err)
	}
}
