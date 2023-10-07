package slog

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
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

func TestLogger_Log(t *testing.T) {
	syncer := &testWriteSyncer{}
	slg := log.NewHelper(NewLogger(slog.New(slog.NewTextHandler(syncer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))))
	slg.Debugw("log", "debug")
	slg.Infow("log", "info")
	slg.Warnw("log", "warn")
	slg.Errorw("log", "error")
	slg.Errorw("log", "error", "except warn")

	except := []string{
		"level=DEBUG msg=\"\" log=debug\n",
		"level=INFO msg=\"\" log=info\n",
		"level=WARN msg=\"\" log=warn\n",
		"level=ERROR msg=\"\" log=error\n",
		"level=WARN msg=\"Keyvalues must appear in pairs: [log error except warn]\"\n",
	}
	for i, s := range except {
		out := strings.Join(strings.Split(syncer.output[i], " ")[1:], " ")
		if s != out {
			t.Logf("except=%s, got=%s", s, out)
			t.Fail()
		}
	}
}
