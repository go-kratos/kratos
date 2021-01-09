package main

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/source/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/log/stdlog"
	servergrpc "github.com/go-kratos/kratos/v2/server/grpc"
	serverhttp "github.com/go-kratos/kratos/v2/server/http"
	transportgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	transporthttp "github.com/go-kratos/kratos/v2/transport/http"

	"google.golang.org/grpc"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "conf", "../../configs", "config path")
}

func main() {
	flag.Parse()
	logger, err := stdlog.NewLogger(stdlog.Writer(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer logger.Close()
	log := log.NewHelper("main", logger)

	c := config.New(config.WithSource(file.NewSource(configPath)))
	if err := c.Load(); err != nil {
		panic(err)
	}
	httpAddr, err := c.Value("server.http.addr").String()
	if err != nil {
		panic(err)
	}
	grpcAddr, err := c.Value("server.grpc.addr").String()
	if err != nil {
		panic(err)
	}
	log.Infof("http listening %s", httpAddr)
	log.Infof("grpc listening %s", grpcAddr)

	httpTransport := transporthttp.NewServer()
	grpcTransport := transportgrpc.NewServer()

	httpServer := serverhttp.NewServer("tcp", httpAddr, serverhttp.ServerHandler(httpTransport))
	grpcServer := servergrpc.NewServer("tcp", grpcAddr, grpc.UnaryInterceptor(grpcTransport.Interceptor()))

	app := kratos.New()
	app.Append(httpServer)
	app.Append(grpcServer)

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		log.Infof("app failed: %v", err)
	}

}
