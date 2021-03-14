package log

import "testing"

func TestStdLogger(t *testing.T) {
	logger := DefaultLogger

	Debug(logger).Print("log", "test debug")
	Info(logger).Print("log", "test info")
	Warn(logger).Print("log", "test warn")
	Error(logger).Print("log", "test error")
}
