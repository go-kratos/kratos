package server

import (
	v1 "github.com/SeeMusic/kratos/examples/i18n/api/helloworld/v1"
	"github.com/SeeMusic/kratos/examples/i18n/internal/conf"
	"github.com/SeeMusic/kratos/examples/i18n/internal/pkg/middleware/localize"
	"github.com/SeeMusic/kratos/examples/i18n/internal/service"
	"github.com/SeeMusic/kratos/v2/log"
	"github.com/SeeMusic/kratos/v2/middleware/logging"
	"github.com/SeeMusic/kratos/v2/middleware/metrics"
	"github.com/SeeMusic/kratos/v2/middleware/recovery"
	"github.com/SeeMusic/kratos/v2/middleware/tracing"
	"github.com/SeeMusic/kratos/v2/middleware/validate"
	"github.com/SeeMusic/kratos/v2/transport/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			metrics.Server(),
			localize.I18N(),
			validate.Validator(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
