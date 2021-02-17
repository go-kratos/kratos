package log

import (
	"testing"
)

func TestHelper(t *testing.T) {
	log := NewHelper("test", DefaultLogger)
	log.Debug("test debug")
	log.Debugf("test %s", "debug")
	log.Debugw("log", "test debug")
}

func TestHelperLevel(t *testing.T) {
	log := NewHelper("test", DefaultLogger)
	log.Debug("test debug")
	log.Info("test info")
	log.Warn("test warn")
	log.Error("test error")
}
