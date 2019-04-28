package http

import (
	"fmt"

	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"
)

// WrapHeader 为返回头添加自定义字段
func WrapHeader(ctx *bm.Context) {
	// Traceid
	tracer, _ := trace.FromContext(ctx.Context)
	traceid := fmt.Sprintf("%s", tracer)
	ctx.Writer.Header().Set("traceid", traceid)

	// Sessionid
	sid := ctx.Request.Header.Get("SessionID")
	if sid == "" {
		sid = SessionID(ctx)
	}
	ctx.Set("SessionID", sid)
	ctx.Writer.Header().Set("SessionID", sid)
}
