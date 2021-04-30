package log

import "testing"

func TestValue(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "caller", DefaultCaller, "ts", DefaultTimestamp)
	logger.Log("msg", "helloworld")

	Debug(logger).Log("msg", "debug value")
	Info(logger).Log("msg", "info value")
	Warn(logger).Log("msg", "warn value")
	Error(logger).Log("msg", "error value")
}
