package logging

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc"
	"time"
)

type loggingServerStream struct {
	req any
	grpc.ServerStream
	logger log.Logger
}

// StreamServer is a server logging middleware for gRPC streams.
func StreamServer(logger log.Logger) middleware.StreamMiddleware {
	return func(handler middleware.StreamHandler) middleware.StreamHandler {
		return func(srv interface{}, stream grpc.ServerStream) error {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)
			ctx := stream.Context()
			startTime := time.Now()
			if info, ok := transport.FromClientContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}
			wrappedStream := &loggingServerStream{
				ServerStream: stream,
				logger:       logger,
			}
			err := handler(srv, wrappedStream)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}
			level, stack := extractError(err)

			log.NewHelper(logger).Log(level,
				"kind", kind,
				"component", kind,
				"operation", operation,
				"args", extractArgs(wrappedStream.req),
				"code", code,
				"reason", reason,
				"stack", stack,
				"latency", time.Since(startTime).Seconds())
			return err
		}
	}
}

func (ss *loggingServerStream) RecvMsg(m interface{}) error {
	var (
		code      int32
		reason    string
		kind      string
		operation string
	)
	err := ss.ServerStream.RecvMsg(m)
	if err != nil {
		level, stack := extractError(err)
		if se := errors.FromError(err); se != nil {
			code = se.Code
			reason = se.Reason
		}
		ctx := ss.Context()
		if info, ok := transport.FromClientContext(ctx); ok {
			kind = info.Kind().String()
			operation = info.Operation()
		}
		log.NewHelper(ss.logger).Log(level,
			"kind", kind,
			"component", kind,
			"operation", operation,
			"code", code,
			"reason", reason,
			"stack", stack,
		)
	}
	if ss.req == nil {
		ss.req = m
	}
	return err
}

func (ss *loggingServerStream) SendMsg(m interface{}) error {
	var (
		code      int32
		reason    string
		kind      string
		operation string
	)
	err := ss.ServerStream.SendMsg(m)
	ctx := ss.Context()
	if info, ok := transport.FromClientContext(ctx); ok {
		kind = info.Kind().String()
		operation = info.Operation()
	}
	if err != nil {
		level, stack := extractError(err)
		if se := errors.FromError(err); se != nil {
			code = se.Code
			reason = se.Reason
		}
		log.NewHelper(ss.logger).Log(level,
			"kind", kind,
			"component", kind,
			"operation", operation,
			"code", code,
			"reason", reason,
			"stack", stack,
		)
	}
	return err
}
