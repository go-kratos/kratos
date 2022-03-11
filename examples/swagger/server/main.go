package main

import (
	"context"
	"fmt"
	"os"

	pb "github.com/SeeMusic/kratos/examples/swagger/helloworld"
	"github.com/SeeMusic/kratos/examples/swagger/reply"
	"github.com/SeeMusic/kratos/v2"
	"github.com/SeeMusic/kratos/v2/errors"
	"github.com/SeeMusic/kratos/v2/log"
	"github.com/SeeMusic/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "helloworld"
	// Version is the version of the compiled software.
	// Version = "v1.0.0"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.Name == "error" {
		return nil, errors.BadRequest("custom_error", fmt.Sprintf("invalid argument %s", in.Name))
	}
	if in.Name == "panic" {
		panic("grpc panic")
	}

	return &pb.HelloReply{Reply: &reply.Reply{Value: fmt.Sprintf("Hello %+v", in.Name)}}, nil
}

func main() {
	stdLogger := log.NewStdLogger(os.Stdout)

	logger := log.NewHelper(stdLogger)
	s := &server{}

	httpSrv := http.NewServer(http.Address(":8000"))
	pb.RegisterGreeterHTTPServer(httpSrv, s)

	h := openapiv2.NewHandler()
	httpSrv.HandlePrefix("/q/", h)

	/*	fs := httpx.FileServer(httpx.Dir("./dist"))
		httpSrv.HandlePrefix("/swaggerui/", httpx.StripPrefix("/swaggerui/", fs))*/

	app := kratos.New(
		kratos.Name(Name),
		kratos.Server(
			httpSrv,
		),
	)

	if err := app.Run(); err != nil {
		logger.Error(err)
	}
}
