package log

import (
	"context"
	"io"
	"os"
	"testing"
)

func TestHelper(_ *testing.T) {
	logger := With(
		DefaultLogger,
		"ts", DefaultTimestamp,
		"caller", DefaultCaller,
		"module", "test",
	)
	log := NewHelper(logger)

	log.Log(LevelDebug, "msg", "test debug")
	log.Debug("test debug")
	log.Debugf("test %s", "debug")
	log.Debugw("log", "test debug")

	log.Warn("test warn")
	log.Warnf("test %s", "warn")
	log.Warnw("log", "test warn")

	subLogger := With(log.Logger(),
		"module", "sub",
	)
	subLog := NewHelper(subLogger)
	subLog.Infof("sub logger test with level %s", "info")
}

func TestHelperWithMsgKey(_ *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	log := NewHelper(logger, WithMessageKey("message"))
	log.Debugf("test %s", "debug")
	log.Debugw("log", "test debug")
}

func TestHelperLevel(_ *testing.T) {
	log := NewHelper(DefaultLogger)
	log.Debug("test debug")
	log.Info("test info")
	log.Infof("test %s", "info")
	log.Warn("test warn")
	log.Error("test error")
	log.Errorf("test %s", "error")
	log.Errorw("log", "test error")
}

func BenchmarkHelperPrint(b *testing.B) {
	log := NewHelper(NewStdLogger(io.Discard))
	for i := 0; i < b.N; i++ {
		log.Debug("test")
	}
}

func BenchmarkHelperPrintFilterLevel(b *testing.B) {
	log := NewHelper(NewFilter(NewStdLogger(io.Discard), FilterLevel(LevelDebug)))
	for i := 0; i < b.N; i++ {
		log.Debug("test")
	}
}

func BenchmarkHelperPrintf(b *testing.B) {
	log := NewHelper(NewStdLogger(io.Discard))
	for i := 0; i < b.N; i++ {
		log.Debugf("%s", "test")
	}
}

func BenchmarkHelperPrintfFilterLevel(b *testing.B) {
	log := NewHelper(NewFilter(NewStdLogger(io.Discard), FilterLevel(LevelInfo)))
	for i := 0; i < b.N; i++ {
		log.Debugf("%s", "test")
	}
}

func BenchmarkHelperPrintw(b *testing.B) {
	log := NewHelper(NewStdLogger(io.Discard))
	for i := 0; i < b.N; i++ {
		log.Debugw("key", "value")
	}
}

type traceKey struct{}

func TestContext(_ *testing.T) {
	logger := With(NewStdLogger(os.Stdout),
		"trace", Trace(),
	)
	log := NewHelper(logger)
	ctx := context.WithValue(context.Background(), traceKey{}, "2233")
	log.WithContext(ctx).Info("got trace!")
}

func Trace() Valuer {
	return func(ctx context.Context) interface{} {
		s, ok := ctx.Value(traceKey{}).(string)
		if !ok {
			return nil
		}
		return s
	}
}
