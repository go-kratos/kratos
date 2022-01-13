package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockErr struct{}

func (*mockErr) Error() string {
	return "mock error"
}

func TestWarp(t *testing.T) {
	var err error = &mockErr{}
	err2 := fmt.Errorf("wrap %w", err)
	assert.Equal(t, err, Unwrap(err2))
	assert.True(t, Is(err2, err))
	err3 := &mockErr{}
	assert.True(t, As(err2, &err3))
}
