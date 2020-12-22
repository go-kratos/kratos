package log

import (
	"testing"
)

type testLogger struct {
	*testing.T
}

func (t *testLogger) Print(kvpiar ...interface{}) {
	t.Log(kvpiar...)
}

func TestLogger(t *testing.T) {
	log := &testLogger{t}
	log.Print("log", "test")

	Debug(log).Print("log", "test")
	Info(log).Print("log", "test")
	Warn(log).Print("log", "test")
	Error(log).Print("log", "test")
}
