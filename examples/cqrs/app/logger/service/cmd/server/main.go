package main

import (
	"flag"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"kratos-cqrs/app/logger/service/internal/conf"
	"kratos-cqrs/pkg/util/bootstrap"
)

// go build -ldflags "-X main.Service.Version=x.y.z"
var (
	Service = bootstrap.NewServiceInfo(
		"kratos.logger.service",
		"1.0.0",
		"",
	)

	Flags = bootstrap.NewCommandFlags()
)

func init() {
	Flags.Init()
}

func newApp(logger log.Logger, gs *grpc.Server, rr registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.ID(Service.GetInstanceId()),
		kratos.Name(Service.Name),
		kratos.Version(Service.Version),
		kratos.Metadata(Service.Metadata),
		kratos.Logger(logger),
		kratos.Server(
			gs,
		),
		kratos.Registrar(rr),
	)
}

func loadConfig() (*conf.Bootstrap, *conf.Registry) {
	c := bootstrap.NewConfigProvider(Flags.ConfigType, Flags.ConfigHost, Flags.Conf, Service.Name)

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	var rc conf.Registry
	if err := c.Scan(&rc); err != nil {
		panic(err)
	}

	return &bc, &rc
}

func main() {
	flag.Parse()

	logger := bootstrap.NewLoggerProvider(&Service)

	bc, rc := loadConfig()
	if bc == nil || rc == nil {
		panic("load config failed")
	}

	err := bootstrap.NewTracerProvider(bc.Trace.Endpoint, Flags.Env, &Service)
	if err != nil {
		panic(err)
	}

	app, cleanup, err := initApp(bc.Server, rc, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
