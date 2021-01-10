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

func (t *testLogger) Close() error {
	return nil
}

func TestLogger(t *testing.T) {
	log := &testLogger{t}

	Debug(log).Print("log", "test debug")
	Info(log).Print("log", "test info")
	Warn(log).Print("log", "test warn")
	Error(log).Print("log", "test error")
}
