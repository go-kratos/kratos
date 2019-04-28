package status

import (
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-common/library/ecode"
	"go-common/library/ecode/pb"
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
func FromError(err error) *status.Status {
	err = errors.Cause(err)
	if code, ok := err.(ecode.Codes); ok {
		// TODO: deal with err
		if gst, err := gRPCStatusFromEcode(code); err == nil {
			return gst
		}
	}
	gst, _ := status.FromError(err)
	return gst
}

func gRPCStatusFromEcode(code ecode.Codes) (*status.Status, error) {
	var st *ecode.Status
	switch v := code.(type) {
	// compatible old pb.Error remove it after nobody use pb.Error.
	case *pb.Error:
		return status.New(codes.Unknown, v.Error()).WithDetails(v)
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
	// gst := status.New(togRPCCode(st), st.Message())
	// NOTE: compatible with PHP swoole gRPC put code in status message as string.
	// gst := status.New(togRPCCode(st), strconv.Itoa(st.Code()))
	gst := status.New(codes.Unknown, strconv.Itoa(st.Code()))
	pbe := &pb.Error{ErrCode: int32(st.Code()), ErrMessage: gst.Message()}
	// NOTE: server return ecode.Status will be covert to pb.Error details will be ignored
	// and put it at details[0] for compatible old client
	return gst.WithDetails(pbe, st.Proto())
}

// ToEcode convert grpc.status to ecode.Codes
func ToEcode(gst *status.Status) ecode.Codes {
	details := gst.Details()
	// reverse range details, details may contain three case,
	// if details contain pb.Error and ecode.Status use eocde.Status first.
	//
	// Details layout:
	// pb.Error [0: pb.Error]
	// both pb.Error and ecode.Status [0: pb.Error, 1: ecode.Status]
	// ecode.Status [0: ecode.Status]
	for i := len(details) - 1; i >= 0; i-- {
		detail := details[i]
		// compatible with old pb.Error.
		if pe, ok := detail.(*pb.Error); ok {
			st := ecode.Error(ecode.Code(pe.ErrCode), pe.ErrMessage)
			if pe.ErrDetail != nil {
				dynMsg := new(ptypes.DynamicAny)
				// TODO deal with unmarshalAny error.
				if err := ptypes.UnmarshalAny(pe.ErrDetail, dynMsg); err == nil {
					st, _ = st.WithDetails(dynMsg.Message)
				}
			}
			return st
		}
		// convert detail to status only use first detail
		if pb, ok := detail.(proto.Message); ok {
			return ecode.FromProto(pb)
		}
	}
	return toECode(gst)
}
