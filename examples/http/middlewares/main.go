package main

import (
	"log"

	"github.com/SeeMusic/kratos/examples/helloworld/helloworld"
	"github.com/SeeMusic/kratos/v2"
	"github.com/SeeMusic/kratos/v2/transport/http"
)

// this example shows how to add middlewares,
// execution order is globalFilter(http) --> routeFilter(http) --> pathFilter(http) --> serviceFilter(service)
func main() {
	s := &server{}
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			// add service filter
			serviceMiddleware,
			serviceMiddleware2,
		),
		// add global filter
		http.Filter(globalFilter, globalFilter2),
	)
	// register http handler to http server
	helloworld.RegisterGreeterHTTPServer(httpSrv, s)

	// add route filter
	r := httpSrv.Route("/", routeFilter, routeFilter2)
	// add path filter to custom route
	r.GET("/hello/{name}", sayHelloHandler, pathFilter, pathFilter2)

	// add path filter to router group
	r2 := r.Group("/v2", pathFilter, pathFilter2)
	r2.GET("/say/{name}", sayHelloHandler)

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
