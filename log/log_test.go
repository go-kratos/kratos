package log

import (
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := DefaultLogger
	Debug(logger).Log("msg", "test debug")
	Info(logger).Log("msg", "test info")
	Warn(logger).Log("msg", "test warn")
	Error(logger).Log("msg", "test error")
}

func TestInfo(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "caller", DefaultCaller, "ts", DefaultTimestamp)
	infoLogger := Info(logger)
	infoLogger.Log("key1", "value1")
	infoLogger.Log("key2", "value2")
	infoLogger.Log("key3", "value3")
}

func TestWrapper(t *testing.T) {
	out := NewStdLogger(os.Stdout)
	err := NewStdLogger(os.Stderr)

	l := With(MultiLogger(out, err), "caller", DefaultCaller, "ts", DefaultTimestamp)
	l.Log("msg", "test")
}
