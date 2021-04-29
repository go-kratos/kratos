package log

import "testing"

func TestValue(t *testing.T) {
	logger := With(DefaultLogger, "caller", Caller(3))
	logger.Print("message", "helloworld")

	Debug(logger).Print("message", "debug value")
	Info(logger).Print("message", "info value")
	Warn(logger).Print("message", "warn value")
	Error(logger).Print("message", "error value")
}
