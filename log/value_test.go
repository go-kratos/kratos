package log

import "testing"

func TestValue(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "caller", DefaultCaller, "ts", DefaultTimestamp)
	logger.Print("msg", "helloworld")

	Debug(logger).Print("msg", "debug value")
	Info(logger).Print("msg", "info value")
	Warn(logger).Print("msg", "warn value")
	Error(logger).Print("msg", "error value")
}
