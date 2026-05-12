package zerolog

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/rs/zerolog"

	klog "github.com/go-kratos/kratos/v2/log"
)

type testWriteSyncer struct {
	output []string
}

func (x *testWriteSyncer) Write(p []byte) (n int, err error) {
	x.output = append(x.output, string(p))
	return len(p), nil
}

func TestLogger(t *testing.T) {
	syncer := &testWriteSyncer{}
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"
	zlogger := zerolog.New(syncer)
	logger := NewLogger(&zlogger)

	logger.Debug("", "log", "debug")
	logger.Info("", "log", "info")
	logger.Warn("", "log", "warn")
	logger.Error("", "log", "error")

	except := []string{
		"{\"level\":\"debug\",\"log\":\"debug\"}\n",
		"{\"level\":\"info\",\"log\":\"info\"}\n",
		"{\"level\":\"warn\",\"log\":\"warn\"}\n",
		"{\"level\":\"error\",\"log\":\"error\"}\n",
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
	zlogger := zerolog.New(syncer)
	logger := NewLogger(&zlogger,
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
