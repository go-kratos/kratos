package logging

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func grpcServerLog(logger log.Logger, ctx context.Context, args string, err error) {
	info, ok := grpc.FromServerContext(ctx)
	if !ok {
		return
	}
	code, errMsg := extractError(err)
	level := log.LevelInfo
	if err != nil {
		level = log.LevelError
	}
	log.WithContext(ctx, logger).Log(level,
		"kind", "server",
		"component", "grpc",
		"grpc.target", info.FullMethod,
		"grpc.args", args,
		"grpc.code", code,
		"grpc.error", errMsg,
	)
}

func grpcClientLog(logger log.Logger, ctx context.Context, args string, err error) {
	info, ok := grpc.FromClientContext(ctx)
	if !ok {
		return
	}
	code, errMsg := extractError(err)
	level := log.LevelInfo
	if err != nil {
		level = log.LevelError
	}
	log.WithContext(ctx, logger).Log(level,
		"kind", "client",
		"component", "grpc",
		"grpc.target", info.FullMethod,
		"grpc.args", args,
		"grpc.code", code,
		"grpc.error", errMsg,
	)
}
