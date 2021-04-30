package log

import "testing"

func TestStdLogger(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "caller", DefaultCaller, "ts", DefaultTimestamp)

	Debug(logger).Print("msg", "test debug")
	Info(logger).Print("msg", "test info")
	Warn(logger).Print("msg", "test warn")
	Error(logger).Print("msg", "test error")
}
