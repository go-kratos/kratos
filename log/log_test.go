package log

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	log := NewStdLogger(os.Stdout)
	log.Print("log", "test")

	Debug(log).Print("log", "test")
	Info(log).Print("log", "test")
	Warn(log).Print("log", "test")
	Error(log).Print("log", "test")
}

func BenchmarkLoggerPrint(b *testing.B) {
	b.SetParallelism(100)
	log := NewStdLogger(ioutil.Discard)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print("log", "test")
		}
	})
}
func BenchmarkLoggerHelperV(b *testing.B) {
	b.SetParallelism(100)
	log := NewHelper("test", NewStdLogger(ioutil.Discard))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.V(10).Print("log", "test")
		}
	})
}
func BenchmarkLoggerHelperInfo(b *testing.B) {
	b.SetParallelism(100)
	log := NewHelper("test", NewStdLogger(ioutil.Discard))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("log", "test")
		}
	})
}
func BenchmarkLoggerHelperInfof(b *testing.B) {
	b.SetParallelism(100)
	log := NewHelper("test", NewStdLogger(ioutil.Discard))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Infof("log %s", "test")
		}
	})
}
func BenchmarkLoggerHelperInfow(b *testing.B) {
	b.SetParallelism(100)
	log := NewHelper("test", NewStdLogger(ioutil.Discard))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Infow("log", "test")
		}
	})
}
