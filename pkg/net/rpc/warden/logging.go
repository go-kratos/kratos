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
)

// Warden Log Flag
const (
	// disable all log.
	LogFlagDisable = 1 << iota
	// disable print args on log.
	LogFlagDisableArgs
	// disable info level log.
	LogFlagDisableInfo
)

type logOption struct {
	grpc.EmptyDialOption
	grpc.EmptyCallOption
	flag int8
}

// WithLogFlag disable client access log.
func WithLogFlag(flag int8) grpc.CallOption {
	return logOption{flag: flag}
}

// WithDialLogFlag set client level log behaviour.
func WithDialLogFlag(flag int8) grpc.DialOption {
	return logOption{flag: flag}
}

func extractLogCallOption(opts []grpc.CallOption) (flag int8) {
	for _, opt := range opts {
		if logOpt, ok := opt.(logOption); ok {
			return logOpt.flag
		}
	}
	return
}

func extractLogDialOption(opts []grpc.DialOption) (flag int8) {
	for _, opt := range opts {
		if logOpt, ok := opt.(logOption); ok {
			return logOpt.flag
		}
	}
	return
}

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
func clientLogging(dialOptions ...grpc.DialOption) grpc.UnaryClientInterceptor {
	defaultFlag := extractLogDialOption(dialOptions)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		logFlag := extractLogCallOption(opts) | defaultFlag

		startTime := time.Now()
		var peerInfo peer.Peer
		opts = append(opts, grpc.Peer(&peerInfo))

		// invoker requests
		err := invoker(ctx, method, req, reply, cc, opts...)

		// after request
		code := ecode.Cause(err).Code()
		duration := time.Since(startTime)
		// monitor
		_metricClientReqDur.Observe(int64(duration/time.Millisecond), method)
		_metricClientReqCodeTotal.Inc(method, strconv.Itoa(code))

		if logFlag&LogFlagDisable != 0 {
			return err
		}
		// TODO: find better way to deal with slow log.
		if logFlag&LogFlagDisableInfo != 0 && err == nil && duration < 500*time.Millisecond {
			return err
		}
		logFields := make([]log.D, 0, 7)
		logFields = append(logFields, log.KVString("path", method))
		logFields = append(logFields, log.KVInt("ret", code))
		logFields = append(logFields, log.KVFloat64("ts", duration.Seconds()))
		logFields = append(logFields, log.KVString("source", "grpc-access-log"))
		if peerInfo.Addr != nil {
			logFields = append(logFields, log.KVString("ip", peerInfo.Addr.String()))
		}
		if logFlag&LogFlagDisableArgs == 0 {
			// TODO: it will panic if someone remove String method from protobuf message struct that auto generate from protoc.
			logFields = append(logFields, log.KVString("args", req.(fmt.Stringer).String()))
		}
		if err != nil {
			logFields = append(logFields, log.KVString("error", err.Error()), log.KVString("stack", fmt.Sprintf("%+v", err)))
		}
		logFn(code, duration)(ctx, logFields...)
		return err
	}
}

// serverLogging warden grpc logging
func serverLogging(logFlag int8) grpc.UnaryServerInterceptor {
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
		_metricServerReqDur.Observe(int64(duration/time.Millisecond), info.FullMethod, caller)
		_metricServerReqCodeTotal.Inc(info.FullMethod, caller, strconv.Itoa(code))

		if logFlag&LogFlagDisable != 0 {
			return resp, err
		}
		// TODO: find better way to deal with slow log.
		if logFlag&LogFlagDisableInfo != 0 && err == nil && duration < 500*time.Millisecond {
			return resp, err
		}
		logFields := []log.D{
			log.KVString("user", caller),
			log.KVString("ip", remoteIP),
			log.KVString("path", info.FullMethod),
			log.KVInt("ret", code),
			log.KVFloat64("ts", duration.Seconds()),
			log.KVFloat64("timeout_quota", quota),
			log.KVString("source", "grpc-access-log"),
		}
		if logFlag&LogFlagDisableArgs == 0 {
			// TODO: it will panic if someone remove String method from protobuf message struct that auto generate from protoc.
			logFields = append(logFields, log.KVString("args", req.(fmt.Stringer).String()))
		}
		if err != nil {
			logFields = append(logFields, log.KVString("error", err.Error()), log.KVString("stack", fmt.Sprintf("%+v", err)))
		}
		logFn(code, duration)(ctx, logFields...)
		return resp, err
	}
}
