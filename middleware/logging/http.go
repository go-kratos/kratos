package logging

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func httpServerLog(logger log.Logger, ctx context.Context, args string, err error) {
	info, ok := http.FromServerContext(ctx)
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
		"component", "http",
		"http.target", info.Request.RequestURI,
		"http.method", info.Request.Method,
		"http.args", args,
		"http.code", code,
		"http.error", errMsg,
	)
}

func httpClientLog(logger log.Logger, ctx context.Context, args string, err error) {
	info, ok := http.FromClientContext(ctx)
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
		"component", "http",
		"http.target", info.Request.RequestURI,
		"http.method", info.Request.Method,
		"http.args", args,
		"http.code", code,
		"http.error", errMsg,
	)
}
