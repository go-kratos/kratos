package stdlog

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestLogger(t *testing.T) {
	logger := NewLogger(Writer(os.Stdout))
	logger.Print("log", "test")

	log.Debug(logger).Print("log", "test")
	log.Info(logger).Print("log", "test")
	log.Warn(logger).Print("log", "test")
	log.Error(logger).Print("log", "test")
}

func BenchmarkLoggerPrint(b *testing.B) {
	b.SetParallelism(100)
	log := NewLogger(Writer(ioutil.Discard))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print("log", "test")
		}
	})
}
func BenchmarkLoggerHelperV(b *testing.B) {
	b.SetParallelism(100)
	log := log.NewHelper("test", NewLogger(Writer(ioutil.Discard)))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.V(10).Print("log", "test")
		}
	})
}
func BenchmarkLoggerHelperInfo(b *testing.B) {
	b.SetParallelism(100)
	log := log.NewHelper("test", NewLogger(Writer(ioutil.Discard)))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("log", "test")
		}
	})
}
func BenchmarkLoggerHelperInfof(b *testing.B) {
	b.SetParallelism(100)
	log := log.NewHelper("test", NewLogger(Writer(ioutil.Discard)))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Infof("log %s", "test")
		}
	})
}
func BenchmarkLoggerHelperInfow(b *testing.B) {
	b.SetParallelism(100)
	log := log.NewHelper("test", NewLogger(Writer(ioutil.Discard)))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Infow("log", "test")
		}
	})
}
