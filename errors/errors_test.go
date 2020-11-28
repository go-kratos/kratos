package errors

import (
	"errors"
	"testing"
)

func TestErrorsMatch(t *testing.T) {
	s := &Error{Code: 1}
	st := &Error{Code: 2}

	if errors.Is(s, st) {
		t.Errorf("error is not match: %+v -> %+v", s, st)
	}

	s.Code = 1
	st.Code = 1
	if !errors.Is(s, st) {
		t.Errorf("error is not match: %+v -> %+v", s, st)
	}

	s.WithDetails(&ErrorInfo{Reason: "test_reason"})
	st.WithDetails(&ErrorInfo{Reason: "test_reason"})

	if !errors.Is(s, st) {
		t.Errorf("error is not match: %+v -> %+v", s, st)
	}
}
