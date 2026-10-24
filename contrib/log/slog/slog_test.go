package slog

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

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
	slg := log.NewHelper(NewLogger(slog.NewTextHandler(syncer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})))
	slg.Debugw("msg", "debug")
	slg.Infow("msg", "info")
	slg.Warnw("msg", "warn")
	slg.Errorw("msg", "error")
	slg.Errorw("msg", "error", "except_warn")

	except := []string{
		"level=DEBUG msg=debug\n",
		"level=INFO msg=info\n",
		"level=WARN msg=warn\n",
		"level=ERROR msg=error\n",
		"level=WARN msg=\"Keyvalues must appear in pairs: [msg error except_warn]\"\n",
	}
	for i, s := range except {
		if s != syncer.output[i] {
			t.Errorf("except=%q, got=%q", s, syncer.output[i])
		}
	}
}

func TestLogger_Log_with(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		AddSource: false,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC))
			}
			return a
		},
	}))

	h := log.NewHelper(log.With(l, "name", "foo"))

	h.Infow("msg", "hello", "user", "bar")

	expect := "time=2006-01-02T15:04:05.000Z level=INFO msg=hello name=foo user=bar\n"
	if got := buf.String(); got != expect {
		t.Fatalf("expect=%q, got=%q,", expect, got)
	}
}

func TestLogger_Log_another_msg_key(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		AddSource: false,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	h := log.NewHelper(l, log.WithMessageKey("other"))
	h.Info("hello")

	expect := "level=INFO msg=\"\" other=hello\n"
	if got := buf.String(); got != expect {
		t.Fatalf("expect=%q, got=%q,", expect, got)
	}
}

func TestLogger_Log_no_msg_key(t *testing.T) {
	var buf bytes.Buffer
	h := log.NewHelper(NewLogger(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		AddSource: false,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})))
	h.Infow("key", "value")

	expect := "level=INFO msg=\"\" key=value\n"
	if got := buf.String(); got != expect {
		t.Fatalf("expect=%q, got=%q,", expect, got)
	}
}
