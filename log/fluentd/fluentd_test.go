package fluentd

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

var opts = []Option{FluentHost("127.0.0.1"), FluentPort(24224)}

func TestLogger(t *testing.T) {
	logger := NewLogger(opts...)
	logger.Print("log", "test")

	log.Debug(logger).Print("log", "test")
	log.Info(logger).Print("log", "test")
	log.Warn(logger).Print("log", "test")
	log.Error(logger).Print("log", "test")
}

func BenchmarkLoggerPrint(b *testing.B) {
	b.SetParallelism(100)
	log := NewLogger(opts...)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print("log", "test")
		}
	})
}
func BenchmarkLoggerHelperV(b *testing.B) {
	b.SetParallelism(100)
	log := log.NewHelper("test", NewLogger(opts...))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.V(10).Print("log", "test")
		}
	})
}
func BenchmarkLoggerHelperInfo(b *testing.B) {
	b.SetParallelism(100)
	log := log.NewHelper("test", NewLogger(opts...))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("log", "test")
		}
	})
}
func BenchmarkLoggerHelperInfof(b *testing.B) {
	b.SetParallelism(100)
	log := log.NewHelper("test", NewLogger(opts...))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Infof("log %s", "test")
		}
	})
}
func BenchmarkLoggerHelperInfow(b *testing.B) {
	b.SetParallelism(100)
	log := log.NewHelper("test", NewLogger(opts...))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Infow("log", "test")
		}
	})
}
