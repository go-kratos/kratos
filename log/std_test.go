package log

import "testing"

func TestStdLogger(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "caller", DefaultCaller, "ts", DefaultTimestamp)

	Debug(logger).Log("msg", "test debug")
	Info(logger).Log("msg", "test info")
	Warn(logger).Log("msg", "test warn")
	Error(logger).Log("msg", "test error")
}
