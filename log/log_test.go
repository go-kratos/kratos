package log

import (
	"bytes"
	"context"
	"sync"
	"testing"
	"time"
)

func TestInfo(_ *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp)
	logger = With(logger, "caller", DefaultCaller)
	_ = logger.Log(LevelInfo, "key1", "value1")
}

type traceIDKey struct{}

func setTraceID(ctx context.Context, tid string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, tid)
}

func traceIDValuer() Valuer {
	return func(ctx context.Context) any {
		if ctx == nil {
			return ""
		}
		return ctx.Value(traceIDKey{})
	}
}

func TestWithContext(t *testing.T) {
	var buf bytes.Buffer
	pctx := context.Background()
	l := NewFilter(
		With(NewStdLogger(&buf), "trace-id", traceIDValuer()),
		FilterLevel(LevelInfo),
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		NewHelper(l).Info("done1")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		tid := "world"
		ctx := setTraceID(pctx, tid)
		NewHelper((WithContext(ctx, l))).Info("done2")
	}()

	wg.Wait()
	expected := "INFO trace-id=world msg=done2\nINFO trace-id= msg=done1\n"
	if got := buf.String(); got != expected {
		t.Errorf("got: %#v", got)
	}
}
