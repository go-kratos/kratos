package log

import (
	"os"
	"testing"
)

func TestWrapper(t *testing.T) {
	out := NewStdLogger(os.Stdout)
	err := NewStdLogger(os.Stderr)

	l := Wrap(out, err)
	l.Print("message", "test")
}
