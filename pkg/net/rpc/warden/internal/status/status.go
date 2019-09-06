package status

import (
	"context"
	"strconv"

	"github.com/bilibili/kratos/pkg/ecode"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// togRPCCode convert ecode.Codo to gRPC code
func togRPCCode(code ecode.Codes) codes.Code {
	switch code.Code() {
	case ecode.OK.Code():
		return codes.OK
	case ecode.RequestErr.Code():
		return codes.InvalidArgument
	case ecode.NothingFound.Code():
		return codes.NotFound
	case ecode.Unauthorized.Code():
		return codes.Unauthenticated
	case ecode.AccessDenied.Code():
		return codes.PermissionDenied
	case ecode.LimitExceed.Code():
		return codes.ResourceExhausted
	case ecode.MethodNotAllowed.Code():
		return codes.Unimplemented
	case ecode.Deadline.Code():
		return codes.DeadlineExceeded
	case ecode.ServiceUnavailable.Code():
		return codes.Unavailable
	}
	return codes.Unknown
}

func toECode(gst *status.Status) ecode.Code {
	gcode := gst.Code()
	switch gcode {
	case codes.OK:
		return ecode.OK
	case codes.InvalidArgument:
		return ecode.RequestErr
	case codes.NotFound:
		return ecode.NothingFound
	case codes.PermissionDenied:
		return ecode.AccessDenied
	case codes.Unauthenticated:
		return ecode.Unauthorized
	case codes.ResourceExhausted:
		return ecode.LimitExceed
	case codes.Unimplemented:
		return ecode.MethodNotAllowed
	case codes.DeadlineExceeded:
		return ecode.Deadline
	case codes.Unavailable:
		return ecode.ServiceUnavailable
	case codes.Unknown:
		return ecode.String(gst.Message())
	}
	return ecode.ServerErr
}

// FromError convert error for service reply and try to convert it to grpc.Status.
func FromError(svrErr error) (gst *status.Status) {
	var err error
	svrErr = errors.Cause(svrErr)
	if code, ok := svrErr.(ecode.Codes); ok {
		// TODO: deal with err
		if gst, err = gRPCStatusFromEcode(code); err == nil {
			return
		}
	}
	// for some special error convert context.Canceled to ecode.Canceled,
	// context.DeadlineExceeded to ecode.DeadlineExceeded only for raw error
	// if err be wrapped will not effect.
	switch svrErr {
	case context.Canceled:
		gst, _ = gRPCStatusFromEcode(ecode.Canceled)
	case context.DeadlineExceeded:
		gst, _ = gRPCStatusFromEcode(ecode.Deadline)
	default:
		gst, _ = status.FromError(svrErr)
	}
	return
}

func gRPCStatusFromEcode(code ecode.Codes) (*status.Status, error) {
	var st *ecode.Status
	switch v := code.(type) {
	case *ecode.Status:
		st = v
	case ecode.Code:
		st = ecode.FromCode(v)
	default:
		st = ecode.Error(ecode.Code(code.Code()), code.Message())
		for _, detail := range code.Details() {
			if msg, ok := detail.(proto.Message); ok {
				st.WithDetails(msg)
			}
		}
	}
	gst := status.New(codes.Unknown, strconv.Itoa(st.Code()))
	return gst.WithDetails(st.Proto())
}

// ToEcode convert grpc.status to ecode.Codes
func ToEcode(gst *status.Status) ecode.Codes {
	details := gst.Details()
	for _, detail := range details {
		// convert detail to status only use first detail
		if pb, ok := detail.(proto.Message); ok {
			return ecode.FromProto(pb)
		}
	}
	return toECode(gst)
}
