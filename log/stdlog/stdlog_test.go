package stdlog

import (
	"os"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

type Discard int

func (d Discard) Write(p []byte) (n int, err error) { return }
func (d Discard) Close() (err error)                { return }

func TestLogger(t *testing.T) {
	logger, err := NewLogger(Writer(os.Stdout))
	if err != nil {
		t.Error(err)
	}
	defer logger.Close()

	log.Debug(logger).Print("log", "test debug")
	log.Info(logger).Print("log", "test info")
	log.Warn(logger).Print("log", "test warn")
	log.Error(logger).Print("log", "test error")
}

func BenchmarkLoggerPrint(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger(Writer(Discard(0)))
	if err != nil {
		b.Error(err)
	}
	defer logger.Close()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print("log", "test")
		}
	})
}

func BenchmarkLoggerHelperV(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger(Writer(Discard(0)))
	if err != nil {
		b.Error(err)
	}
	defer logger.Close()
	h := log.NewHelper("test", logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.V(10).Print("log", "test")
		}
	})
}

func BenchmarkLoggerHelperInfo(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger(Writer(Discard(0)))
	if err != nil {
		b.Error(err)
	}
	defer logger.Close()
	h := log.NewHelper("test", logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Info("test")
		}
	})
}

func BenchmarkLoggerHelperInfof(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger(Writer(Discard(0)))
	if err != nil {
		b.Error(err)
	}
	defer logger.Close()
	h := log.NewHelper("test", logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Infof("log %s", "test")
		}
	})
}

func BenchmarkLoggerHelperInfow(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger(Writer(Discard(0)))
	if err != nil {
		b.Error(err)
	}
	defer logger.Close()
	log := log.NewHelper("test", logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Infow("log", "test")
		}
	})
}
