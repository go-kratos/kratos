package log

import "testing"

func TestFilter(t *testing.T) {
	logger := NewFilter(DefaultLogger, LevelInfo)
	logger = With(logger, "caller", DefaultCaller, "ts", DefaultTimestamp)

	Debug(logger).Print("msg", "debug value")
	Info(logger).Print("msg", "info value")
	Warn(logger).Print("msg", "warn value")
	Error(logger).Print("msg", "error value")
}
