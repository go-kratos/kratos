package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestError(t *testing.T) {
	var (
		base *Error
	)
	err := Newf(http.StatusBadRequest, "reason", "message")
	err2 := Newf(http.StatusBadRequest, "reason", "message")
	err3 := err.WithMetadata(map[string]string{
		"foo": "bar",
	})
	werr := fmt.Errorf("wrap %w", err)

	if errors.Is(err, new(Error)) {
		t.Errorf("should not be equal: %v", err)
	}
	if !errors.Is(werr, err) {
		t.Errorf("should be equal: %v", err)
	}
	if !errors.Is(werr, err2) {
		t.Errorf("should be equal: %v", err)
	}

	if !errors.As(err, &base) {
		t.Errorf("should be matchs: %v", err)
	}
	if !IsBadRequest(err) {
		t.Errorf("should be matchs: %v", err)
	}

	if reason := Reason(err); reason != err3.Reason {
		t.Errorf("got %s want: %s", reason, err)
	}

	if err3.Metadata["foo"] != "bar" {
		t.Error("not expected metadata")
	}

	gs := err.GRPCStatus()
	se := FromError(gs.Err())
	if se.Reason != se.Reason {
		t.Errorf("got %+v want %+v", se, err)
	}
}
