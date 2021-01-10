package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/errors"
	pb "github.com/go-kratos/kratos/v2/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/log/stdlog"
	"github.com/go-kratos/kratos/v2/middleware"
	servergrpc "github.com/go-kratos/kratos/v2/server/grpc"
	serverhttp "github.com/go-kratos/kratos/v2/server/http"
	"github.com/go-kratos/kratos/v2/transport"
	transportgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transporthttp "github.com/go-kratos/kratos/v2/transport/http"

	"google.golang.org/grpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.Name == "error" {
		return nil, errors.InvalidArgument("BadRequest", "invalid argument %s", in.Name)
	}
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %+v", in)}, nil
}

func logger1(logger log.Logger) middleware.Middleware {
	log := log.NewHelper("logger1", logger)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			log.Info("before")

			return handler(ctx, req)
		}
	}
}

func logger2(logger log.Logger) middleware.Middleware {
	log := log.NewHelper("logger2", logger)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			resp, err := handler(ctx, req)

			log.Info("after")

			return resp, err
		}
	}
}

func logger3(logger log.Logger) middleware.Middleware {
	log := log.NewHelper("logger2", logger)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			tr, ok := transport.FromContext(ctx)
			if ok {
				log.Infof("transport: %+v", tr)
			}
			h, ok := transporthttp.FromContext(ctx)
			if ok {
				log.Infof("http: [%s] %s", h.Request.Method, h.Request.URL.Path)
			}
			g, ok := transportgrpc.FromContext(ctx)
			if ok {
				log.Infof("grpc: %s", g.FullMethod)
			}

			return handler(ctx, req)
		}
	}
}

func main() {
	logger, err := stdlog.NewLogger(stdlog.Writer(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	log := log.NewHelper("main", logger)

	s := &server{}
	app := kratos.New()

	httpTransport := transporthttp.NewServer(transporthttp.ServerMiddleware(logger1(logger), logger2(logger)))
	httpTransport.Use(s, logger3(logger))

	grpcTransport := transportgrpc.NewServer(transportgrpc.ServerMiddleware(logger1(logger), logger2(logger)))
	grpcTransport.Use(s, logger3(logger))

	httpServer := serverhttp.NewServer("tcp", ":8000", serverhttp.Handler(httpTransport))
	grpcServer := servergrpc.NewServer("tcp", ":9000", grpc.UnaryInterceptor(grpcTransport.Interceptor()))

	pb.RegisterGreeterServer(grpcServer, s)
	pb.RegisterGreeterHTTPServer(httpTransport, s)

	app.Append(httpServer)
	app.Append(grpcServer)

	if err := app.Run(); err != nil {
		log.Error(err)
	}
}
