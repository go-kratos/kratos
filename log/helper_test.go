package log

import (
	"testing"
)

func TestHelper(t *testing.T) {
	log := NewHelper("test", &testLogger{t})
	log.Debug("test debug")
	log.Debugf("test %s", "debug")
	log.Debugw("log", "test debug")
}

func TestHelperLevel(t *testing.T) {
	log := NewHelper("test", &testLogger{t})
	log.Debug("test debug")
	log.Info("test info")
	log.Warn("test warn")
	log.Error("test error")
}

func TestHelperVerbose(t *testing.T) {
	log := NewHelper("test", &testLogger{t})
	log.V(1).Print(LevelDebug, "log", "test debug")
	log.V(5).Print(LevelInfo, "log", "test info")
	log.V(10).Print(LevelWarn, "log", "test warn")
	log.V(15).Print(LevelError, "log", "test error")
}
