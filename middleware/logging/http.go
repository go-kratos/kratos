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
	traceID, spanID := extractTrace(ctx)
	code, errMsg := extractError(err)
	level := log.LevelInfo
	if err != nil {
		level = log.LevelError
	}
	logger.Log(level,
		"kind", "server",
		"component", "http",
		"trace_id", traceID,
		"span_id", spanID,
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
	traceID, spanID := extractTrace(ctx)
	code, errMsg := extractError(err)
	level := log.LevelInfo
	if err != nil {
		level = log.LevelError
	}
	logger.Log(level,
		"kind", "client",
		"component", "http",
		"trace_id", traceID,
		"span_id", spanID,
		"http.target", info.Request.RequestURI,
		"http.method", info.Request.Method,
		"http.args", args,
		"http.code", code,
		"http.error", errMsg,
	)
}
