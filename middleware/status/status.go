package status

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/internal/http"
	"github.com/go-kratos/kratos/v2/middleware"

	//lint:ignore SA1019 grpc
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

// HandlerFunc is middleware error handler.
type HandlerFunc func(context.Context, error) error

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
		handler: encodeErr,
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			reply, err := handler(ctx, req)
			if err != nil {
				return nil, options.handler(ctx, err)
			}
			return reply, nil
		}
	}
}

// Client is an error middleware.
func Client(opts ...Option) middleware.Middleware {
	options := options{
		handler: decodeErr,
	}
	for _, o := range opts {
		o(&options)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			reply, err := handler(ctx, req)
			if err != nil {
				return nil, options.handler(ctx, err)
			}
			return reply, nil
		}
	}
}

func encodeErr(ctx context.Context, err error) error {
	var details []proto.Message
	if target := new(errors.ErrorInfo); errors.As(err, &target) {
		details = append(details, &errdetails.ErrorInfo{
			Domain:   target.Domain,
			Reason:   target.Reason,
			Metadata: target.Metadata,
		})
	}
	es := errors.FromError(err)
	gs := status.New(http.GRPCCodeFromStatus(es.Code), es.Message)
	gs, err = gs.WithDetails(details...)
	if err != nil {
		return err
	}
	return gs.Err()
}

func decodeErr(ctx context.Context, err error) error {
	gs := status.Convert(err)
	code := http.StatusFromGRPCCode(gs.Code())
	message := gs.Message()
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			return errors.Errorf(
				code,
				d.Domain,
				d.Reason,
				message,
			).WithMetadata(d.Metadata)
		}
	}
	return errors.New(code, message)
}
