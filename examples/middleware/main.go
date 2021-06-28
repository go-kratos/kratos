package main

import (
	"context"
	"log"

	"github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type server struct {
	helloworld.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "hello"}, nil
}

func main() {
	s := &server{}
	/*
		The processing order of a request is the order of middleware registration,
		while the processing order of the response is the reverse of the registration order.
		         ┌───────────────────┐
		         │MIDDLEWARE 1       │
		         │ ┌────────────────┐│
		         │ │MIDDLEWARE 2    ││
		         │ │ ┌─────────────┐││
		         │ │ │MIDDLEWARE 3 │││
		         │ │ │ ┌─────────┐ │││
		REQUEST  │ │ │ │  YOUR   │ │││  RESPONSE
		   ──────┼─┼─┼─▷ HANDLER ○─┼┼┼───▷
		         │ │ │ └─────────┘ │││
		         │ │ └─────────────┘││
		         │ └────────────────┘│
		         └───────────────────┘
	*/
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			serviceMiddleware,
			serviceMiddleware2,
		),
	)
	helloworld.RegisterGreeterHTTPServer(httpSrv, s)
	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			httpSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
