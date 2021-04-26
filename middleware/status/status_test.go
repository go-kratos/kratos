package status

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/errors"
)

func TestErrEncoder(t *testing.T) {
	err := errors.BadRequest("test", "invalid_argument", "format")
	en := encodeErr(context.Background(), err)
	de := decodeErr(context.Background(), en)
	if !errors.IsBadRequest(de) {
		t.Errorf("expected %v got %v", err, de)
	}
}
