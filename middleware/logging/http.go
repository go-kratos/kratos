package logging

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func httpServerLog(logger log.Logger, ctx context.Context, args string, err error) {
	var level = log.LevelInfo
	var code = 200
	var errMsg = ""

	info, ok := http.FromServerContext(ctx)
	if !ok {
		return
	}

	if err != nil {
		level = log.LevelError
		code, errMsg = extractError(err)
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
	var level = log.LevelInfo
	var code = 200
	var errMsg = ""

	info, ok := http.FromClientContext(ctx)
	if !ok {
		return
	}

	if err != nil {
		level = log.LevelError
		code, errMsg = extractError(err)
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
