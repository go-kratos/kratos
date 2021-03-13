package log

import (
	"io"
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

func BenchmarkHelperPrint(b *testing.B) {
	log := NewHelper("test", NewStdLogger(io.Discard))
	for i := 0; i < b.N; i++ {
		log.Debug("test")
	}
}

func BenchmarkHelperPrintf(b *testing.B) {
	log := NewHelper("test", NewStdLogger(io.Discard))
	for i := 0; i < b.N; i++ {
		log.Debugf("%s", "test")
	}
}

func BenchmarkHelperPrintw(b *testing.B) {
	log := NewHelper("test", NewStdLogger(io.Discard))
	for i := 0; i < b.N; i++ {
		log.Debugw("key", "value")
	}
}
