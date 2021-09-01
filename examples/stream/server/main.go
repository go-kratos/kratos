package main

import (
	"log"

	"github.com/go-kratos/kratos/examples/stream/hello"
	"github.com/go-kratos/kratos/examples/stream/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func main() {
	grpcSrv := grpc.NewServer(
		grpc.Address(":9001"),
		grpc.Middleware(
			recovery.Recovery(),
		),
	)
	hello.RegisterHelloServer(grpcSrv, service.NewHelloService())

	app := kratos.New(
		kratos.Name("hello"),
		kratos.Server(
			grpcSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
