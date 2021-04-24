package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	var (
		base *Error
		br   *BadRequestError
	)
	err := NewBadRequest("domain", "reason", "message")
	werr := fmt.Errorf("wrap %w", err)

	if errors.Is(err, base) {
		t.Errorf("should not be equal: %v", err)
	}
	if !errors.Is(werr, err) {
		t.Errorf("should be equal: %v", err)
	}

	if !errors.As(err, &base) {
		t.Errorf("should be matchs: %v", err)
	}
	if !errors.As(err, &br) {
		t.Errorf("should be matchs: %v", err)
	}
}
