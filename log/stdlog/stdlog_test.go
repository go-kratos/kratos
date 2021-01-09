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
	logger.Print("log", "test")

	log.Debug(logger).Print("log", "test")
	log.Info(logger).Print("log", "test")
	log.Warn(logger).Print("log", "test")
	log.Error(logger).Print("log", "test")
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
	log := log.NewHelper("test", logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.V(10).Print("log", "test")
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
	log := log.NewHelper("test", logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("log", "test")
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
	log := log.NewHelper("test", logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Infof("log %s", "test")
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
