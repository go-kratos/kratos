package warden

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"github.com/bilibili/kratos/pkg/ecode"
	"github.com/bilibili/kratos/pkg/log"
	"github.com/bilibili/kratos/pkg/net/metadata"
	"github.com/bilibili/kratos/pkg/stat"
)

var (
	statsClient = stat.RPCClient
	statsServer = stat.RPCServer
)

func logFn(code int, dt time.Duration) func(context.Context, ...log.D) {
	switch {
	case code < 0:
		return log.Errorv
	case dt >= time.Millisecond*500:
		// TODO: slowlog make it configurable.
		return log.Warnv
	case code > 0:
		return log.Warnv
	}
	return log.Infov
}

// clientLogging warden grpc logging
func clientLogging() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		startTime := time.Now()
		var peerInfo peer.Peer
		opts = append(opts, grpc.Peer(&peerInfo))

		// invoker requests
		err := invoker(ctx, method, req, reply, cc, opts...)

		// after request
		code := ecode.Cause(err).Code()
		duration := time.Since(startTime)
		// monitor
		statsClient.Timing(method, int64(duration/time.Millisecond))
		statsClient.Incr(method, strconv.Itoa(code))

		var ip string
		if peerInfo.Addr != nil {
			ip = peerInfo.Addr.String()
		}
		logFields := []log.D{
			log.KVString("ip", ip),
			log.KVString("path", method),
			log.KVInt("ret", code),
			// TODO: it will panic if someone remove String method from protobuf message struct that auto generate from protoc.
			log.KVString("args", req.(fmt.Stringer).String()),
			log.KVFloat64("ts", duration.Seconds()),
			log.KVString("source", "grpc-access-log"),
		}
		if err != nil {
			logFields = append(logFields, log.KV("error", err.Error()), log.KVString("stack", fmt.Sprintf("%+v", err)))
		}
		logFn(code, duration)(ctx, logFields...)
		return err
	}
}

// serverLogging warden grpc logging
func serverLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		caller := metadata.String(ctx, metadata.Caller)
		if caller == "" {
			caller = "no_user"
		}
		var remoteIP string
		if peerInfo, ok := peer.FromContext(ctx); ok {
			remoteIP = peerInfo.Addr.String()
		}
		var quota float64
		if deadline, ok := ctx.Deadline(); ok {
			quota = time.Until(deadline).Seconds()
		}

		// call server handler
		resp, err := handler(ctx, req)

		// after server response
		code := ecode.Cause(err).Code()
		duration := time.Since(startTime)

		// monitor
		statsServer.Timing(caller, int64(duration/time.Millisecond), info.FullMethod)
		statsServer.Incr(caller, info.FullMethod, strconv.Itoa(code))
		logFields := []log.D{
			log.KVString("user", caller),
			log.KVString("ip", remoteIP),
			log.KVString("path", info.FullMethod),
			log.KVInt("ret", code),
			// TODO: it will panic if someone remove String method from protobuf message struct that auto generate from protoc.
			log.KVString("args", req.(fmt.Stringer).String()),
			log.KVFloat64("ts", duration.Seconds()),
			log.KVFloat64("timeout_quota", quota),
			log.KVString("source", "grpc-access-log"),
		}
		if err != nil {
			logFields = append(logFields, log.KV("error", err.Error()), log.KV("stack", fmt.Sprintf("%+v", err)))
		}
		logFn(code, duration)(ctx, logFields...)
		return resp, err
	}
}
