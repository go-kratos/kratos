package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorsMatch(t *testing.T) {
	s := &StatusError{Code: 1}
	st := &StatusError{Code: 2}

	if errors.Is(s, st) {
		t.Errorf("error is not match: %+v -> %+v", s, st)
	}

	s.Code = 1
	st.Code = 1
	if !errors.Is(s, st) {
		t.Errorf("error is not match: %+v -> %+v", s, st)
	}

	s.Reason = "test_reason"
	s.Reason = "test_reason"

	if !errors.Is(s, st) {
		t.Errorf("error is not match: %+v -> %+v", s, st)
	}

	if Reason(s) != "test_reason" {
		t.Errorf("error is not match: %+v -> %+v", s, st)
	}
}

func TestErrorIs(t *testing.T) {
	err1 := &StatusError{Code: 1}
	t.Log(err1)
	err2 := fmt.Errorf("wrap : %w", err1)
	t.Log(err2)

	if !(errors.Is(err2, err1)) {
		t.Errorf("error is not match: a: %v b: %v ", err2, err1)
	}
}

func TestErrorAs(t *testing.T) {
	err1 := &StatusError{Code: 1}
	err2 := fmt.Errorf("wrap : %w", err1)

	err3 := new(StatusError)
	if !errors.As(err2, &err3) {
		t.Errorf("error is not match: %v", err2)
	}
}
