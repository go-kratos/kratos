package log

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"strconv"
	"time"

	"github.com/bilibili/kratos/pkg/conf/env"
	"github.com/bilibili/kratos/pkg/log/internal/core"
	"github.com/bilibili/kratos/pkg/net/metadata"
	"github.com/bilibili/kratos/pkg/net/trace"
)

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
	fields[_appID] = c.AppID
	fields[_instanceID] = c.Host
	if metadata.Bool(ctx, metadata.Mirror) {
		fields[_mirror] = true
	}
}

// toMap convert D slice to map[string]interface{} for legacy file and stdout.
func toMap(args ...D) map[string]interface{} {
	d := make(map[string]interface{}, 10+len(args))
	for _, arg := range args {
		switch arg.Type {
		case core.UintType, core.Uint64Type, core.IntTpye, core.Int64Type:
			d[arg.Key] = arg.Int64Val
		case core.StringType:
			d[arg.Key] = arg.StringVal
		case core.Float32Type:
			d[arg.Key] = math.Float32frombits(uint32(arg.Int64Val))
		case core.Float64Type:
			d[arg.Key] = math.Float64frombits(uint64(arg.Int64Val))
		case core.DurationType:
			d[arg.Key] = time.Duration(arg.Int64Val)
		default:
			d[arg.Key] = arg.Value
		}
	}
	return d
}

// funcName get func name.
func funcName(skip int) (name string) {
	if _, file, lineNo, ok := runtime.Caller(skip); ok {
		return file + ":" + strconv.Itoa(lineNo)
	}
	return "unknown:0"
}
