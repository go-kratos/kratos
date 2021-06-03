package log

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
)

func TestHelper(t *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	log := NewHelper(logger)

	log.Log(LevelDebug, "msg", "test debug")
	log.Debug("test debug")
	log.Debugf("test %s", "debug")
	log.Debugw("log", "test debug")
}

func TestHelperLevel(t *testing.T) {
	log := NewHelper(DefaultLogger)
	log.Debug("test debug")
	log.Info("test info")
	log.Warn("test warn")
	log.Error("test error")
}

func BenchmarkHelperPrint(b *testing.B) {
	log := NewHelper(NewStdLogger(ioutil.Discard))
	for i := 0; i < b.N; i++ {
		log.Debug("test")
	}
}

func BenchmarkHelperPrintf(b *testing.B) {
	log := NewHelper(NewStdLogger(ioutil.Discard))
	for i := 0; i < b.N; i++ {
		log.Debugf("%s", "test")
	}
}

func BenchmarkHelperPrintw(b *testing.B) {
	log := NewHelper(NewStdLogger(ioutil.Discard))
	for i := 0; i < b.N; i++ {
		log.Debugw("key", "value")
	}
}

func TestContext(t *testing.T) {
	logger := With(NewStdLogger(os.Stdout),
		"trace", Trace(),
	)
	log := NewHelper(logger)
	ctx := context.WithValue(context.Background(), "trace_id", "2233")
	log.WithContext(ctx).Info("got trace!")
}

func Trace() Valuer {
	return func(ctx context.Context) interface{} {
		s := ctx.Value("trace_id").(string)
		return s
	}
}
