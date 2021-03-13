package log

import (
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := DefaultLogger
	Debug(logger).Print("log", "test debug")
	Info(logger).Print("log", "test info")
	Warn(logger).Print("log", "test warn")
	Error(logger).Print("log", "test error")
}

func TestWrapper(t *testing.T) {
	out := NewStdLogger(os.Stdout)
	err := NewStdLogger(os.Stderr)

	l := Wrap(out, err)
	l.Print("message", "test")
}
