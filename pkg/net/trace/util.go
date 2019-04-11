package trace

import (
	"context"
	"encoding/binary"
	"math/rand"
	"time"

	"github.com/bilibili/kratos/pkg/conf/env"
	"github.com/bilibili/kratos/pkg/net/ip"

	"github.com/pkg/errors"
)

var _hostHash byte

func init() {
	rand.Seed(time.Now().UnixNano())
	_hostHash = byte(oneAtTimeHash(env.Hostname))
}

func extendTag() (tags []Tag) {
	tags = append(tags,
		TagString("region", env.Region),
		TagString("zone", env.Zone),
		TagString("hostname", env.Hostname),
		TagString("ip", ip.InternalIP()),
	)
	return
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

var _ctxkey ctxKey = "kratos/pkg/net/trace.trace"

// FromContext returns the trace bound to the context, if any.
func FromContext(ctx context.Context) (t Trace, ok bool) {
	t, ok = ctx.Value(_ctxkey).(Trace)
	return
}

// NewContext new a trace context.
// NOTE: This method is not thread safe.
func NewContext(ctx context.Context, t Trace) context.Context {
	return context.WithValue(ctx, _ctxkey, t)
}
