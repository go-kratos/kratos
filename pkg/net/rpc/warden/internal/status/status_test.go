package status

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	pkgerr "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bilibili/kratos/pkg/ecode"
)

func TestCodeConvert(t *testing.T) {
	var table = map[codes.Code]ecode.Code{
		codes.OK: ecode.OK,
		// codes.Canceled
		codes.Unknown:          ecode.ServerErr,
		codes.InvalidArgument:  ecode.RequestErr,
		codes.DeadlineExceeded: ecode.Deadline,
		codes.NotFound:         ecode.NothingFound,
		// codes.AlreadyExists
		codes.PermissionDenied:  ecode.AccessDenied,
		codes.ResourceExhausted: ecode.LimitExceed,
		// codes.FailedPrecondition
		// codes.Aborted
		// codes.OutOfRange
		codes.Unimplemented: ecode.MethodNotAllowed,
		codes.Unavailable:   ecode.ServiceUnavailable,
		// codes.DataLoss
		codes.Unauthenticated: ecode.Unauthorized,
	}
	for k, v := range table {
		assert.Equal(t, toECode(status.New(k, "-500")), v)
	}
	for k, v := range table {
		assert.Equal(t, togRPCCode(v), k, fmt.Sprintf("togRPC code error: %d -> %d", v, k))
	}
}

func TestNoDetailsConvert(t *testing.T) {
	gst := status.New(codes.Unknown, "-2233")
	assert.Equal(t, toECode(gst).Code(), -2233)

	gst = status.New(codes.Internal, "")
	assert.Equal(t, toECode(gst).Code(), -500)
}

func TestFromError(t *testing.T) {
	t.Run("input general error", func(t *testing.T) {
		err := errors.New("general error")
		gst := FromError(err)

		assert.Equal(t, codes.Unknown, gst.Code())
		assert.Contains(t, gst.Message(), "general")
	})
	t.Run("input wrap error", func(t *testing.T) {
		err := pkgerr.Wrap(ecode.RequestErr, "hh")
		gst := FromError(err)

		assert.Equal(t, "-400", gst.Message())
	})
	t.Run("input ecode.Code", func(t *testing.T) {
		err := ecode.RequestErr
		gst := FromError(err)

		//assert.Equal(t, codes.InvalidArgument, gst.Code())
		// NOTE: set all grpc.status as Unkown when error is ecode.Codes for compatible
		assert.Equal(t, codes.Unknown, gst.Code())
		// NOTE: gst.Message == str(ecode.Code) for compatible php leagcy code
		assert.Equal(t, err.Message(), gst.Message())
	})
	t.Run("input raw Canceled", func(t *testing.T) {
		gst := FromError(context.Canceled)

		assert.Equal(t, codes.Unknown, gst.Code())
		assert.Equal(t, "-498", gst.Message())
	})
	t.Run("input raw DeadlineExceeded", func(t *testing.T) {
		gst := FromError(context.DeadlineExceeded)

		assert.Equal(t, codes.Unknown, gst.Code())
		assert.Equal(t, "-504", gst.Message())
	})
	t.Run("input ecode.Status", func(t *testing.T) {
		m := &timestamp.Timestamp{Seconds: time.Now().Unix()}
		err, _ := ecode.Error(ecode.Unauthorized, "unauthorized").WithDetails(m)
		gst := FromError(err)

		//assert.Equal(t, codes.Unauthenticated, gst.Code())
		// NOTE: set all grpc.status as Unkown when error is ecode.Codes for compatible
		assert.Equal(t, codes.Unknown, gst.Code())
		assert.Len(t, gst.Details(), 1)
		details := gst.Details()
		assert.IsType(t, err.Proto(), details[0])
	})
}

func TestToEcode(t *testing.T) {
	t.Run("input general grpc.Status", func(t *testing.T) {
		gst := status.New(codes.Unknown, "unknown")
		ec := ToEcode(gst)

		assert.Equal(t, int(ecode.ServerErr), ec.Code())
		assert.Equal(t, "-500", ec.Message())
		assert.Len(t, ec.Details(), 0)
	})
	t.Run("input encode.Status", func(t *testing.T) {
		m := &timestamp.Timestamp{Seconds: time.Now().Unix()}
		st, _ := ecode.Errorf(ecode.Unauthorized, "Unauthorized").WithDetails(m)
		gst := status.New(codes.InvalidArgument, "requesterr")
		gst, _ = gst.WithDetails(st.Proto())
		ec := ToEcode(gst)

		assert.Equal(t, int(ecode.Unauthorized), ec.Code())
		assert.Equal(t, "Unauthorized", ec.Message())
		assert.Len(t, ec.Details(), 1)
		assert.IsType(t, m, ec.Details()[0])
	})
}
