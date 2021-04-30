package log

import (
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	logger := DefaultLogger
	Debug(logger).Print("msg", "test debug")
	Info(logger).Print("msg", "test info")
	Warn(logger).Print("msg", "test warn")
	Error(logger).Print("msg", "test error")
}

func TestInfo(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "caller", DefaultCaller, "ts", DefaultTimestamp)
	infoLogger := Info(logger)
	infoLogger.Print("key1", "value1")
	infoLogger.Print("key2", "value2")
	infoLogger.Print("key3", "value3")
}

func TestWrapper(t *testing.T) {
	out := NewStdLogger(os.Stdout)
	err := NewStdLogger(os.Stderr)

	l := With(MultiLogger(out, err), "caller", DefaultCaller, "ts", DefaultTimestamp)
	l.Print("msg", "test")
}
