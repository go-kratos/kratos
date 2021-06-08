package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	transhttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	m := middleware.Chain(
		recovery.Recovery(),
		//...
	)
	router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		//can use other lib such as gin.binding : Bind(r, &req)
		next := func(ctx context.Context, req interface{}) (interface{}, error) {
			//call XXX_service.Func
			fmt.Fprint(w, "Hello Gorilla Mux!")
			return nil, nil
		}
		m(next)
		_, _ = next(r.Context(), &r)
		return
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
		log.Println(err)
	}
}
