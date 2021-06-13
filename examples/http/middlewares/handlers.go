package main

import (
	"context"

	"github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func sayHelloHandler(ctx http.Context) error {
	var in helloworld.HelloRequest
	if err := ctx.Bind(&in); err != nil {
		return err
	}

	// binding /hello/{name} to in.Name
	if err := ctx.BindVars(&in); err != nil {
		return err
	}

	http.SetOperation(ctx, "/helloworld.Greeter/SayHello")
	h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
		return &helloworld.HelloReply{Message: "test:" + req.(*helloworld.HelloRequest).Name}, nil
	})
	return ctx.Returns(h(ctx, &in))
}
