package main

import (
	"context"
	"fmt"

	pb "github.com/go-kratos/kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/binding"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.Name == "error" {
		return nil, errors.BadRequest("custom_error", fmt.Sprintf("invalid argument %s", in.Name))
	}
	if in.Name == "panic" {
		panic("grpc panic")
	}
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %+v", in.Name)}, nil
}

func testHandler(ctx http.Context) error {
	var in pb.HelloRequest
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
