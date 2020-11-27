package status

import (
	"errors"
	"testing"
)

func TestStatusMatch(t *testing.T) {
	s := &Status{Code: 1}
	st := &Status{Code: 2}

	if errors.Is(s, st) {
		t.Errorf("error is not match: %+v -> %+v", s, st)
	}

	s.Code = 1
	st.Code = 1
	if !errors.Is(s, st) {
		t.Errorf("error is not match: %+v -> %+v", s, st)
	}

	s.WithDetails(&ErrorInfo{Reason: "test_reason", Domain: "test_domain"})
	st.WithDetails(&ErrorInfo{Reason: "test_reason", Domain: "test_domain"})

	if !errors.Is(s, st) {
		t.Errorf("error is not match: %+v -> %+v", s, st)
	}
}
