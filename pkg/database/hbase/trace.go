package hbase

import (
	"context"
	"io"

	"github.com/tsuna/gohbase/hrpc"

	"github.com/bilibili/kratos/pkg/net/trace"
)

// TraceHook create new hbase trace hook.
func TraceHook(component, instance string) HookFunc {
	var internalTags []trace.Tag
	internalTags = append(internalTags, trace.TagString(trace.TagComponent, component))
	internalTags = append(internalTags, trace.TagString(trace.TagDBInstance, instance))
	internalTags = append(internalTags, trace.TagString(trace.TagPeerService, "hbase"))
	internalTags = append(internalTags, trace.TagString(trace.TagSpanKind, "client"))
	return func(ctx context.Context, call hrpc.Call, customName string) func(err error) {
		noop := func(error) {}
		root, ok := trace.FromContext(ctx)
		if !ok {
			return noop
		}
		if customName == "" {
			customName = call.Name()
		}
		span := root.Fork("", "Hbase:"+customName)
		span.SetTag(internalTags...)
		statement := string(call.Table()) + " " + string(call.Key())
		span.SetTag(trace.TagString(trace.TagDBStatement, statement))
		return func(err error) {
			if err == io.EOF {
				// reset error for trace.
				err = nil
			}
			span.Finish(&err)
		}
	}
}
