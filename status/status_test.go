package status

import (
	"errors"
	"testing"
)

func TestStatusMatch(t *testing.T) {
	st := &Status{}
	st.WithDetails(&ErrorInfo{Reason: "test_reason", Domain: "test_domain"})

	err := &ErrorInfo{Reason: "test_reason", Domain: "test_domain"}
	if !errors.Is(st, err) {
		t.Errorf("error is not match: %+v", err)
	}
}
