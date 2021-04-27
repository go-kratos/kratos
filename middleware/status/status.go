package status

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"

	//lint:ignore SA1019 grpc
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
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
	gs := status.New(httpToGRPCCode(es.Code), es.Message)
	gs, err = gs.WithDetails(details...)
	if err != nil {
		return err
	}
	return gs.Err()
}

func decodeErr(ctx context.Context, err error) error {
	gs := status.Convert(err)
	code := grpcToHTTPCode(gs.Code())
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

func httpToGRPCCode(code int) codes.Code {
	switch code {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.Aborted
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	}
	return codes.Unknown
}

func grpcToHTTPCode(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Aborted:
		return http.StatusConflict
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	}
	return http.StatusInternalServerError
}
