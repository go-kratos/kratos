package zap

import (
	"log/slog"
	"strings"
	"testing"

	klog "github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type testWriteSyncer struct {
	output []string
}

func (x *testWriteSyncer) Write(p []byte) (n int, err error) {
	x.output = append(x.output, string(p))
	return len(p), nil
}

func (x *testWriteSyncer) Sync() error {
	return nil
}

func TestLogger(t *testing.T) {
	syncer := &testWriteSyncer{}
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), syncer, zap.DebugLevel)
	zlogger := zap.New(core).WithOptions()
	handler := NewHandler(zlogger).(*Handler)
	logger := slog.New(handler)

	defer func() { _ = handler.Close() }()

	logger.Debug("", "log", "debug")
	logger.Info("", "log", "info")
	logger.Warn("", "log", "warn")
	logger.Error("", "log", "error")
	logger.Info("hello world")

	except := []string{
		"{\"level\":\"debug\",\"msg\":\"\",\"log\":\"debug\"}\n",
		"{\"level\":\"info\",\"msg\":\"\",\"log\":\"info\"}\n",
		"{\"level\":\"warn\",\"msg\":\"\",\"log\":\"warn\"}\n",
		"{\"level\":\"error\",\"msg\":\"\",\"log\":\"error\"}\n",
		"{\"level\":\"info\",\"msg\":\"hello world\"}\n", // not {"level":"info","msg":"","msg":"hello world"}
	}
	for i, s := range except {
		if s != syncer.output[i] {
			t.Logf("except=%s, got=%s", s, syncer.output[i])
			t.Fail()
		}
	}
}

func TestNewLoggerWithOptions(t *testing.T) {
	syncer := &testWriteSyncer{}
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), syncer, zap.DebugLevel)
	zlogger := zap.New(core)
	logger := NewLogger(zlogger,
		klog.WithAttrs(slog.String("service.name", "helloworld")),
		klog.WithFilter(klog.FilterKey("password")),
	)

	logger.Info("login", "password", "secret")

	if len(syncer.output) != 1 {
		t.Fatalf("len(output) = %d, want 1", len(syncer.output))
	}
	got := syncer.output[0]
	for _, want := range []string{`"service.name":"helloworld"`, `"password":"***"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("output %q does not contain %s", got, want)
		}
	}
}
