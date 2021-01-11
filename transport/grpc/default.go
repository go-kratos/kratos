package grpc

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DefaultErrorEncoder is default error encoder.
func DefaultErrorEncoder(ctx context.Context, err error) error {
	se, ok := err.(*errors.StatusError)
	if !ok {
		se = &errors.StatusError{
			Code:    2,
			Message: "Unknown: " + err.Error(),
		}
	}
	gs := status.Newf(codes.Code(se.Code), "%s: %s", se.Reason, se.Message)
	gs, err = gs.WithDetails(&errdetails.ErrorInfo{
		Reason:   se.Reason,
		Metadata: map[string]string{"message": se.Message},
	})
	if err != nil {
		return err
	}
	return gs.Err()
}

// DefaultErrorDecoder is default error decoder.
func DefaultErrorDecoder(ctx context.Context, err error) error {
	gs := status.Convert(err)
	var (
		reason  string
		message string
	)
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			reason = d.Reason
			message = d.Metadata["message"]
		}
	}
	return &errors.StatusError{
		Code:    int32(gs.Code()),
		Reason:  reason,
		Message: message,
		Details: gs.Details(),
	}
}
