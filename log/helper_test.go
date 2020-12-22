package log

import (
	"testing"
)

func TestHelper(t *testing.T) {
	log := NewHelper("test", &testLogger{t})
	log.Debug("test log")
	log.Debugf("test %s", "log")
	log.Debugw("test", "log")
}

func TestHelperLevel(t *testing.T) {
	log := NewHelper("test", &testLogger{t})
	log.Debug("test log")
	log.Info("test log")
	log.Warn("test log")
	log.Error("test log")
}

func TestHelperVerbose(t *testing.T) {
	log := NewHelper("test", &testLogger{t})
	log.V(1).Print("log", "test log")
	log.V(5).Print("log", "test log")
	log.V(10).Print("log", "test log")
	log.V(15).Print("log", "test log")
}
