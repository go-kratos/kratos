package kratos

import (
	"log"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func ExampleApp() {
	// init & defer
	// TODO do something

	// transport server
	httpSrv := http.NewServer(http.WithAddress(":8000"))
	grpcSrv := grpc.NewServer(grpc.WithAddress(":9000"))

	// application lifecycle
	app := New()
	app.Append(Hook{OnStart: httpSrv.Start, OnStop: httpSrv.Stop})
	app.Append(Hook{OnStart: grpcSrv.Start, OnStop: grpcSrv.Stop})

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		log.Printf("app failed: %v\n", err)
	}
}
