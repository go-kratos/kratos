package errors

import (
	"fmt"
	"testing"
)

type mockErr struct{}

func (*mockErr) Error() string {
	return "mock error"
}

func TestWarp(t *testing.T) {
	var err error = &mockErr{}
	err2 := fmt.Errorf("wrap %w", err)
	if err != Unwrap(err2) {
		t.Errorf("got %v want: %v", err, Unwrap(err2))
	}
	if !Is(err2, err) {
		t.Errorf("Is(err2, err) got %v want: %v", Is(err2, err), true)
	}
	err3 := &mockErr{}
	if !As(err2, &err3) {
		t.Errorf("As(err2, &err3) got %v want: %v", As(err2, &err3), true)
	}
}
