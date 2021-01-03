package kratos

import (
	"log"

	servergrpc "github.com/go-kratos/kratos/v2/server/grpc"
	serverhttp "github.com/go-kratos/kratos/v2/server/http"
	transportgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transporthttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/grpc"
)

func ExampleApp() {

	httpTransport := transporthttp.NewServer()
	grpcTransport := transportgrpc.NewServer()

	// transport server
	httpServer := serverhttp.NewServer("tcp", ":8000", serverhttp.ServerHandler(httpTransport))
	grpcServer := servergrpc.NewServer("tcp", ":9000", grpc.UnaryInterceptor(grpcTransport.Interceptor()))

	// application lifecycle
	app := New()
	app.Append(httpServer)
	app.Append(grpcServer)

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		log.Printf("app failed: %v\n", err)
	}

}
