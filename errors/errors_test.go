package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	var (
		base *Error
	)
	err := BadRequest("translate", "API_KEY_INVALID", "API key not valid. Please pass a valid API key.")
	werr := fmt.Errorf("wrap %w", err)

	if errors.Is(err, new(Error)) {
		t.Errorf("should not be equal: %v", err)
	}
	if !errors.Is(werr, err) {
		t.Errorf("should be equal: %v", err)
	}
	if !errors.Is(werr, BadRequest("translate", "API_KEY_INVALID", "")) {
		t.Errorf("should be equal: %v", err)
	}

	if !errors.As(err, &base) {
		t.Errorf("should be matchs: %v", err)
	}
	if !IsBadRequest(err) {
		t.Errorf("should be matchs: %v", err)
	}

	if reason := Reason(err); reason != "API_KEY_INVALID" {
		t.Errorf("got %s want: %s", reason, err)
	}
}
