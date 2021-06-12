package main

import (
	"context"

	"github.com/go-kratos/kratos/examples/helloworld/helloworld"
	pb "github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
)

func sayHelloHandler(ctx http.Context) error {
	var in helloworld.HelloRequest
	if err := ctx.Bind(&in); err != nil {
		return err
	}

	if err := binding.BindVars(ctx.Vars(), &in); err != nil {
		return err
	}

	transport.SetOperation(ctx, "/helloworld.Greeter/SayHello")
	h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
		return &pb.HelloReply{Message: "test:" + req.(*pb.HelloRequest).Name}, nil
	})
	out, err := h(ctx, &in)
	if err != nil {
		return err
	}
	reply := out.(*pb.HelloReply)
	return ctx.Result(200, reply)
}
