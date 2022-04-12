package log

import (
	"bytes"
	"context"
	"io"
	"testing"
)

func TestValue(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	_ = logger.Log(LevelInfo, "msg", "helloworld")

	logger = DefaultLogger
	logger = With(logger)
	_ = logger.Log(LevelDebug, "msg", "helloworld")

	var v1 interface{}
	got := Value(context.Background(), v1)
	if got != v1 {
		t.Errorf("Value() = %v, want %v", got, v1)
	}
	var v2 Valuer = func(ctx context.Context) interface{} {
		return 3
	}
	got = Value(context.Background(), v2)
	res := got.(int)
	if res != 3 {
		t.Errorf("Value() = %v, want %v", res, 3)
	}
}

func TestCaller(t *testing.T) {
	var buf bytes.Buffer
	h := NewHelper(With(NewStdLogger(&buf), "caller", DefaultCaller))
	h.Infow("args", "test")
	want := "INFO caller=value_test.go:37,log.TestCaller args=test\n"
	if want != buf.String() {
		t.Fatalf("want %s，have %s", want, buf.String())
	}
}

func TestCallerLine(t *testing.T) {
	var buf bytes.Buffer
	h := NewHelper(With(NewStdLogger(&buf), "callerLine", DefaultCallerLine))
	h.Infow("args", "test")
	want := "INFO callerLine=47 args=test\n"
	if want != buf.String() {
		t.Fatalf("want %s，have %s", want, buf.String())
	}
}

func TestCallerFile(t *testing.T) {
	var buf bytes.Buffer
	h := NewHelper(With(NewStdLogger(&buf), "callerFile", DefaultCallerFile))
	h.Infow("args", "test")
	want := "INFO callerFile=value_test.go args=test\n"
	if want != buf.String() {
		t.Fatalf("want %s，have %s", want, buf.String())
	}
}

func TestCallerFuncName(t *testing.T) {
	var buf bytes.Buffer
	h := NewHelper(With(NewStdLogger(&buf), "callerFuncName", DefaultCallerFuncName))
	h.Infow("args", "test")
	want := "INFO callerFuncName=log.TestCallerFuncName args=test\n"
	if want != buf.String() {
		t.Fatalf("want %s，have %s", want, buf.String())
	}
}

func BenchmarkCaller(b *testing.B) {
	h := NewHelper(With(NewStdLogger(io.Discard), "caller", DefaultCaller))
	for i := 0; i < b.N; i++ {
		h.Infow("args", "test")
	}
}

func BenchmarkCallerFile(b *testing.B) {
	h := NewHelper(With(NewStdLogger(io.Discard), "callerFile", DefaultCallerFile))
	for i := 0; i < b.N; i++ {
		h.Infow("args", "test")
	}
}

func BenchmarkCallerFuncName(b *testing.B) {
	h := NewHelper(With(NewStdLogger(io.Discard), "callerFuncName", DefaultCallerFuncName))
	for i := 0; i < b.N; i++ {
		h.Infow("args", "test")
	}
}

func BenchmarkCallerLine(b *testing.B) {
	h := NewHelper(With(NewStdLogger(io.Discard), "callerLine", DefaultCallerLine))
	for i := 0; i < b.N; i++ {
		h.Infow("args", "test")
	}
}
