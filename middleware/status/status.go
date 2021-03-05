package status

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"

	//lint:ignore SA1019 grpc
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HandlerFunc is middleware error handler.
type HandlerFunc func(error) error

// Option is recovery option.
type Option func(*options)

type options struct {
	handler HandlerFunc
}

// WithHandler with status handler.
func WithHandler(h HandlerFunc) Option {
	return func(o *options) {
		o.handler = h
	}
}

// Server is an error middleware.
func Server(opts ...Option) middleware.Middleware {
	options := options{
		handler: errorEncode,
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			reply, err := handler(ctx, req)
			if err != nil {
				return nil, options.handler(err)
			}
			return reply, nil
		}
	}
}

// Client is an error middleware.
func Client(opts ...Option) middleware.Middleware {
	options := options{
		handler: errorDecode,
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			reply, err := handler(ctx, req)
			if err != nil {
				return nil, options.handler(err)
			}
			return reply, nil
		}
	}
}

func errorEncode(err error) error {
	se, ok := errors.FromError(err)
	if !ok {
		se = &errors.StatusError{
			Code:    2,
			Message: err.Error(),
		}
	}
	gs := status.Newf(codes.Code(se.Code), "%s: %s", se.Reason, se.Message)
	details := []proto.Message{
		&errdetails.ErrorInfo{
			Reason:   se.Reason,
			Metadata: map[string]string{"message": se.Message},
		},
	}
	for _, any := range se.Details {
		detail := &ptypes.DynamicAny{}
		if err := ptypes.UnmarshalAny(any, detail); err != nil {
			continue
		}
		details = append(details, detail.Message)
	}
	gs, err = gs.WithDetails(details...)
	if err != nil {
		return err
	}
	return gs.Err()
}

func errorDecode(err error) error {
	gs := status.Convert(err)
	se := &errors.StatusError{
		Code:    int32(gs.Code()),
		Details: gs.Proto().Details,
	}
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			se.Reason = d.Reason
			se.Message = d.Metadata["message"]
			return se
		}
	}
	return se
}
