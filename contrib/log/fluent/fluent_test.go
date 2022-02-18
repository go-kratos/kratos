package fluent

import (
	"io"
	"net"
	"os"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestMain(m *testing.M) {
	if ln, err := net.Listen("tcp", ":24224"); err == nil {
		defer ln.Close()
		go func() {
			for {
				conn, err := ln.Accept()
				if err != nil {
					return
				}
				defer conn.Close()
				if _, err = io.ReadAll(conn); err != nil {
					continue
				}
			}
		}()
	}

	os.Exit(m.Run())
}

func TestLogger(t *testing.T) {
	logger, err := NewLogger("tcp://127.0.0.1:24224")
	if err != nil {
		t.Error(err)
	}
	defer logger.Close()
	flog := log.NewHelper(logger)

	flog.Debug("log", "test")
	flog.Info("log", "test")
	flog.Warn("log", "test")
	flog.Error("log", "test")
}

func BenchmarkLoggerPrint(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger("tcp://127.0.0.1:24224")
	flog := log.NewHelper(logger)
	if err != nil {
		b.Error(err)
	}
	defer logger.Close()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			flog.Info("log", "test")
		}
	})
}

func BenchmarkLoggerHelperV(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger("tcp://127.0.0.1:24224")
	if err != nil {
		b.Error(err)
	}
	h := log.NewHelper(logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Info("log", "test")
		}
	})
}

func BenchmarkLoggerHelperInfo(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger("tcp://127.0.0.1:24224")
	if err != nil {
		b.Error(err)
	}
	defer logger.Close()
	h := log.NewHelper(logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Info("test")
		}
	})
}

func BenchmarkLoggerHelperInfof(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger("tcp://127.0.0.1:24224")
	if err != nil {
		b.Error(err)
	}
	defer logger.Close()
	h := log.NewHelper(logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Infof("log %s", "test")
		}
	})
}

func BenchmarkLoggerHelperInfow(b *testing.B) {
	b.SetParallelism(100)
	logger, err := NewLogger("tcp://127.0.0.1:24224")
	if err != nil {
		b.Error(err)
	}
	defer logger.Close()
	h := log.NewHelper(logger)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.Infow("log", "test")
		}
	})
}
