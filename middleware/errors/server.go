package errors

import (
	"context"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

const (
	ServiceName = "sn"
)

func metadata(err error) map[string]string {
	if e, ok := err.(*errors.Error); ok {
		if e.Metadata == nil {
			e.Metadata = make(map[string]string)
		}
		return e.Metadata
	}
	gs, ok := status.FromError(err)
	if ok {
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				if d.Metadata == nil {
					d.Metadata = make(map[string]string)
				}
				return d.Metadata
			}
		}
	}
	return nil
}

// Server is an server errors middleware.
func Server() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			reply, err = handler(ctx, req)
			if err != nil {
				md := metadata(err)
				if md != nil {
					if info, ok := kratos.FromContext(ctx); ok {
						if md[ServiceName] != "" && md[ServiceName] != info.Name() {
							err = errors.New(errors.UnknownCode, errors.UnknownReason, "")
						} else {
							md[ServiceName] = info.Name()
						}
					}
				}
			}
			return reply, err
		}
	}
}
