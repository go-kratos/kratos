package status

import (
	"testing"

	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrEncoder(t *testing.T) {
	err := errors.BadRequest("InvalidArgument", "format")
	en := encodeErr(err)
	if code := status.Code(en); code != codes.InvalidArgument {
		t.Errorf("expected %d got %d", codes.InvalidArgument, code)
	}
	de := decodeErr(en)
	if !errors.IsBadRequest(de) {
		t.Errorf("expected %v got %v", err, de)
	}
}
