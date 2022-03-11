package server

import (
	v1 "github.com/SeeMusic/kratos/examples/transaction/api/transaction/v1"
	"github.com/SeeMusic/kratos/examples/transaction/gorm/internal/conf"
	"github.com/SeeMusic/kratos/examples/transaction/gorm/internal/service"
	"github.com/SeeMusic/kratos/v2/log"
	"github.com/SeeMusic/kratos/v2/middleware/logging"
	"github.com/SeeMusic/kratos/v2/middleware/recovery"
	"github.com/SeeMusic/kratos/v2/middleware/tracing"
	"github.com/SeeMusic/kratos/v2/middleware/validate"
	"github.com/SeeMusic/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, logger log.Logger, transaction *service.TransactionService) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			logging.Server(logger),
			validate.Validator(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterTransactionServiceServer(srv, transaction)
	return srv
}
