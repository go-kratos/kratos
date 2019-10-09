package trace

import (
	"context"
	"encoding/binary"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/bilibili/kratos/pkg/conf/env"
)

var _sdkLogger *log.Logger

var _hostHash byte

func init() {
	rand.Seed(time.Now().UnixNano())
	_hostHash = byte(oneAtTimeHash(env.Hostname))
	_sdkLogger = log.New(os.Stderr, "dapper-client ", log.LstdFlags)
}

func extendTag() (tags []Tag) {
	clientUUID := uuid.New().String()
	tags = append(tags,
		TagString("hostname", env.Hostname),
		TagString("zone", env.Zone),
		TagString("region", env.Region),
		TagString("client-uuid", clientUUID),
	)
	return
}

func serviceNameFromEnv() string {
	return env.AppID
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

const _ctxkey ctxKey = "go-common/net/trace.trace"

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
