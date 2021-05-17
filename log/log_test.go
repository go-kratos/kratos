package log

import (
	"os"
	"testing"
)

func TestInfo(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	logger.Log(LevelInfo, "key1", "value1")
}

func TestWrapper(t *testing.T) {
	out := NewStdLogger(os.Stdout)
	err := NewStdLogger(os.Stderr)

	l := With(MultiLogger(out, err), "ts", DefaultTimestamp, "caller", DefaultCaller)
	l.Log(LevelInfo, "msg", "test")
}
