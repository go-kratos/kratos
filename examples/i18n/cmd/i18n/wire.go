// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos/examples/i18n/internal/biz"
	"github.com/go-kratos/kratos/examples/i18n/internal/conf"
	"github.com/go-kratos/kratos/examples/i18n/internal/data"
	"github.com/go-kratos/kratos/examples/i18n/internal/server"
	"github.com/go-kratos/kratos/examples/i18n/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
