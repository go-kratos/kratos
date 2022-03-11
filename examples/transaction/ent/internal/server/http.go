package server

import (
	v1 "github.com/SeeMusic/kratos/examples/transaction/api/transaction/v1"
	"github.com/SeeMusic/kratos/examples/transaction/ent/internal/conf"
	"github.com/SeeMusic/kratos/examples/transaction/ent/internal/service"
	"github.com/SeeMusic/kratos/v2/log"
	"github.com/SeeMusic/kratos/v2/middleware/logging"
	"github.com/SeeMusic/kratos/v2/middleware/recovery"
	"github.com/SeeMusic/kratos/v2/middleware/tracing"
	"github.com/SeeMusic/kratos/v2/middleware/validate"
	"github.com/SeeMusic/kratos/v2/transport/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, logger log.Logger, transaction *service.TransactionService) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
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
	v1.RegisterTransactionServiceHTTPServer(srv, transaction)
	return srv
}
