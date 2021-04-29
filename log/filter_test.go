package log

import "testing"

func TestFilter(t *testing.T) {
	logger := NewFilter(DefaultLogger, LevelInfo)

	Debug(logger).Print("message", "debug value")
	Info(logger).Print("message", "info value")
	Warn(logger).Print("message", "warn value")
	Error(logger).Print("message", "error value")
}
