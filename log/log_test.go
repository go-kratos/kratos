package log

import (
	"testing"
)

type testLogger struct {
	*testing.T
}

func (t *testLogger) Print(level Level, kvpiar ...interface{}) {
	t.Log(kvpiar...)
}

func (t *testLogger) Close() error {
	return nil
}

func TestLogger(t *testing.T) {
	log := &testLogger{t}

	log.Print(LevelDebug, "log", "test debug")
	log.Print(LevelInfo, "log", "test info")
	log.Print(LevelWarn, "log", "test warn")
	log.Print(LevelError, "log", "test error")
}
