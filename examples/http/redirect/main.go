package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/SeeMusic/kratos/examples/helloworld/helloworld"
	"github.com/SeeMusic/kratos/v2"
	"github.com/SeeMusic/kratos/v2/errors"
	khttp "github.com/SeeMusic/kratos/v2/transport/http"
)

type server struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	if in.Name == "error" {
		return nil, errors.BadRequest("custom_error", fmt.Sprintf("invalid argument %s", in.Name))
	}
	if in.Name == "panic" {
		panic("server panic")
	}
	return &helloworld.HelloReply{Message: fmt.Sprintf("Hello %+v", in.Name)}, nil
}

func redirectFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/helloworld/kratos" {
			http.Redirect(w, r, "https://go-kratos.dev/", http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	httpSrv := khttp.NewServer(
		khttp.Address(":8000"),
		khttp.Filter(redirectFilter),
	)
	s := &server{}
	helloworld.RegisterGreeterHTTPServer(httpSrv, s)

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

// get http://127.0.0.1:8000/helloworld/kratos
// get http://127.0.0.1:8000/helloworld/go-kratos.dev
