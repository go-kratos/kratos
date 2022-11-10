package fluent

import (
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

func TestMain(m *testing.M) {
	listener := func(ln net.Listener) {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_, err = io.ReadAll(conn)
		if err != nil {
			return
		}
	}

	if ln, err := net.Listen("tcp", ":24224"); err == nil {
		defer ln.Close()
		go func() {
			for {
				listener(ln)
			}
		}()
	}

	os.Exit(m.Run())
}

func TestWithTimeout(t *testing.T) {
	opts := new(options)
	var duration time.Duration = 1000000000
	funcTimeout := WithTimeout(duration)
	funcTimeout(opts)
	if opts.timeout != duration {
		t.Errorf("WithTimeout() = %v, want %v", opts.timeout, duration)
	}
}

func TestWithWriteTimeout(t *testing.T) {
	opts := new(options)
	var duration time.Duration = 1000000000
	funcWriteTimeout := WithWriteTimeout(duration)
	funcWriteTimeout(opts)
	if opts.writeTimeout != duration {
		t.Errorf("WithWriteTimeout() = %v, want %v", opts.writeTimeout, duration)
	}
}

func TestWithBufferLimit(t *testing.T) {
	opts := new(options)
	bufferLimit := 1000000000
	funcBufferLimit := WithBufferLimit(bufferLimit)
	funcBufferLimit(opts)
	if opts.bufferLimit != bufferLimit {
		t.Errorf("WithBufferLimit() = %d, want %d", opts.bufferLimit, bufferLimit)
	}
}

func TestWithRetryWait(t *testing.T) {
	opts := new(options)
	retryWait := 1000000000
	funcRetryWait := WithRetryWait(retryWait)
	funcRetryWait(opts)
	if opts.retryWait != retryWait {
		t.Errorf("WithRetryWait() = %d, want %d", opts.retryWait, retryWait)
	}
}

func TestWithMaxRetry(t *testing.T) {
	opts := new(options)
	maxRetry := 1000000000
	funcMaxRetry := WithMaxRetry(maxRetry)
	funcMaxRetry(opts)
	if opts.maxRetry != maxRetry {
		t.Errorf("WithMaxRetry() = %d, want %d", opts.maxRetry, maxRetry)
	}
}

func TestWithMaxRetryWait(t *testing.T) {
	opts := new(options)
	maxRetryWait := 1000000000
	funcMaxRetryWait := WithMaxRetryWait(maxRetryWait)
	funcMaxRetryWait(opts)
	if opts.maxRetryWait != maxRetryWait {
		t.Errorf("WithMaxRetryWait() = %d, want %d", opts.maxRetryWait, maxRetryWait)
	}
}

func TestWithTagPrefix(t *testing.T) {
	opts := new(options)
	tagPrefix := "tag_prefix"
	funcTagPrefix := WithTagPrefix(tagPrefix)
	funcTagPrefix(opts)
	if opts.tagPrefix != tagPrefix {
		t.Errorf("WithTagPrefix() = %s, want %s", opts.tagPrefix, tagPrefix)
	}
}

func TestWithAsync(t *testing.T) {
	opts := new(options)
	async := true
	funcAsync := WithAsync(async)
	funcAsync(opts)
	if opts.async != async {
		t.Errorf("WithAsync() = %t, want %t", opts.async, async)
	}
}

func TestWithForceStopAsyncSend(t *testing.T) {
	opts := new(options)
	forceStopAsyncSend := true
	funcForceStopAsyncSend := WithForceStopAsyncSend(forceStopAsyncSend)
	funcForceStopAsyncSend(opts)
	if opts.forceStopAsyncSend != forceStopAsyncSend {
		t.Errorf("WithForceStopAsyncSend() = %t, want %t", opts.forceStopAsyncSend, forceStopAsyncSend)
	}
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

func TestLoggerWithOpt(t *testing.T) {
	var duration time.Duration = 1000000000
	logger, err := NewLogger("tcp://127.0.0.1:24224", WithTimeout(duration))
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

func TestLoggerError(t *testing.T) {
	errCase := []string{
		"foo",
		"tcp://127.0.0.1/",
		"tcp://127.0.0.1:1234a",
		"tcp://127.0.0.1:65535",
		"https://127.0.0.1:8080",
		"unix://foo/bar",
	}
	for _, errc := range errCase {
		_, err := NewLogger(errc)
		if err == nil {
			t.Error(err)
		}
	}
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
