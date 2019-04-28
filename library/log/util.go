package log

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"sync"

	"go-common/library/conf/env"
	"go-common/library/net/metadata"
	"go-common/library/net/trace"
)

var fm sync.Map

func addExtraField(ctx context.Context, fields map[string]interface{}) {
	if t, ok := trace.FromContext(ctx); ok {
		if s, ok := t.(fmt.Stringer); ok {
			fields[_tid] = s.String()
		} else {
			fields[_tid] = fmt.Sprintf("%s", t)
		}
	}
	if caller := metadata.String(ctx, metadata.Caller); caller != "" {
		fields[_caller] = caller
	}
	if color := metadata.String(ctx, metadata.Color); color != "" {
		fields[_color] = color
	}
	if cluster := metadata.String(ctx, metadata.Cluster); cluster != "" {
		fields[_cluster] = cluster
	}
	fields[_deplyEnv] = env.DeployEnv
	fields[_zone] = env.Zone
	fields[_appID] = c.Family
	fields[_instanceID] = c.Host
	if metadata.Bool(ctx, metadata.Mirror) {
		fields[_mirror] = true
	}
}

// funcName get func name.
func funcName(skip int) (name string) {
	if pc, _, lineNo, ok := runtime.Caller(skip); ok {
		if v, ok := fm.Load(pc); ok {
			name = v.(string)
		} else {
			name = runtime.FuncForPC(pc).Name() + ":" + strconv.FormatInt(int64(lineNo), 10)
			fm.Store(pc, name)
		}
	}
	return
}
