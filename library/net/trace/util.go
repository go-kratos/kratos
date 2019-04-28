package trace

import (
	"context"
	"encoding/binary"
	"math/rand"
	"time"

	"github.com/pkg/errors"

	"go-common/library/conf/env"
	"go-common/library/net/metadata"
)

var _hostHash byte

func init() {
	rand.Seed(time.Now().UnixNano())
	_hostHash = byte(oneAtTimeHash(env.Hostname))
}

func extendTag() (tags []Tag) {
	tags = append(tags,
		TagString("hostname", env.Hostname),
		TagString("ip", env.IP),
		TagString("zone", env.Zone),
		TagString("region", env.Region),
	)
	return
}

func serviceNameFromEnv() string {
	return env.AppID
}

func isUATEnv() bool {
	return env.DeployEnv == env.DeployEnvUat
}

func genID() uint64 {
	var b [8]byte
	// i think this code will not survive to 2106-02-07
	binary.BigEndian.PutUint32(b[4:], uint32(time.Now().Unix())>>8)
	b[4] = _hostHash
	binary.BigEndian.PutUint32(b[:4], uint32(rand.Int31()))
	return binary.BigEndian.Uint64(b[:])
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type ctxKey string

var _ctxkey ctxKey = "go-common/net/trace.trace"

// FromContext returns the trace bound to the context, if any.
func FromContext(ctx context.Context) (t Trace, ok bool) {
	if v := metadata.Value(ctx, metadata.Trace); v != nil {
		t, ok = v.(Trace)
		return
	}
	t, ok = ctx.Value(_ctxkey).(Trace)
	return
}

// NewContext new a trace context.
// NOTE: This method is not thread safe.
func NewContext(ctx context.Context, t Trace) context.Context {
	if md, ok := metadata.FromContext(ctx); ok {
		md[metadata.Trace] = t
		return ctx
	}
	return context.WithValue(ctx, _ctxkey, t)
}
