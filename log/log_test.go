package log

import (
	"context"
	"os"
	"testing"
)

func TestInfo(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	_ = logger.Log(LevelInfo, "key1", "value1")
}

func TestWrapper(t *testing.T) {
	out := NewStdLogger(os.Stdout)
	err := NewStdLogger(os.Stderr)

	l := With(MultiLogger(out, err), "ts", DefaultTimestamp, "caller", DefaultCaller)
	_ = l.Log(LevelInfo, "msg", "test")
}

func TestWithContext(t *testing.T) {
	WithContext(context.Background(), nil)
}
