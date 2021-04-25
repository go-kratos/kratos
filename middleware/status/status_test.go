package status

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrEncoder(t *testing.T) {
	err := errors.BadRequest("InvalidArgument", "format")
	en := encodeErr(context.Background(), err)
	if code := status.Code(en); code != codes.InvalidArgument {
		t.Errorf("expected %d got %d", codes.InvalidArgument, code)
	}
	de := decodeErr(context.Background(), en)
	if !errors.IsBadRequest(de) {
		t.Errorf("expected %v got %v", err, de)
	}
}
