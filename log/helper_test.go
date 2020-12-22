package log

import (
	"os"
	"testing"
)

func TestHelper(t *testing.T) {
	log := GetHelper("", NewStdLogger(os.Stdout))
	log.Debug("test log")
	log.Debugf("test %s", "log")
	log.Debugw("test", "log")
}

func TestHelperLevel(t *testing.T) {
	log := GetHelper("", NewStdLogger(os.Stdout), AllowLevel(LevelInfo))
	log.Debug("test log")
	log.Info("test log")
	log.Warn("test log")
	log.Error("test log")
}

func TestHelperVerbose(t *testing.T) {
	log := GetHelper("", NewStdLogger(os.Stdout), AllowVerbose(10))
	log.V(1).Print("test log")
	log.V(5).Print("test log")
	log.V(10).Print("test log")
	log.V(15).Print("test log")
}
