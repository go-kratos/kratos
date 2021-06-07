package logging

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/internal/httputil"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// grpcServerLog is a server handler when transport is KindGRPC
func grpcServerLog(logger log.Logger, ctx context.Context, args string, err error) {
	info, ok := grpc.FromServerContext(ctx)
	if !ok {
		return
	}

	level := log.LevelInfo
	if err != nil {
		level = log.LevelError
	}

	log.WithContext(ctx, logger).Log(level,
		"kind", "server",
		"component", "grpc",
		"grpc.target", info.FullMethod,
		"grpc.args", args,
		"grpc.code", httputil.GRPCCodeFromStatus(errors.Code(err)),
		"grpc.error", extractError(err),
	)
}

// grpcClientLog is a client handler when transport is KindGRPC
func grpcClientLog(logger log.Logger, ctx context.Context, args string, err error) {
	info, ok := grpc.FromClientContext(ctx)
	if !ok {
		return
	}
	level := log.LevelInfo
	if err != nil {
		level = log.LevelError
	}
	log.WithContext(ctx, logger).Log(level,
		"kind", "client",
		"component", "grpc",
		"grpc.target", info.FullMethod,
		"grpc.args", args,
		"grpc.code", httputil.GRPCCodeFromStatus(errors.Code(err)),
		"grpc.error", extractError(err),
	)
}
