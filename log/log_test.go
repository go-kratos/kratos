package log

import (
	"context"
	"testing"
	"time"
)

func TestInfo(t *testing.T) {
	l := DefaultLogger
	l = With(l)
	l = With(l)
	l = With(l, "unpaired_key")
	c, ok := l.(*logger)
	if !ok || len(c.prefix) != 2 {
		t.Error("test failed")
	}

	l = With(l, "ts", DefaultTimestamp, "caller", DefaultCaller, "unpaired_key")
	l = With(l, "ts", Timestamp(time.ANSIC), "caller", Caller(-1), "unpaired_key")
	c, ok = l.(*logger)
	if !ok || len(c.prefix) != 14 {
		t.Error("test failed")
	}

	_ = l.Log(LevelInfo, "key1", "value1")
}

func TestWithContext(t *testing.T) {
	WithContext(context.Background(), nil)
}

func TestWithReplace(t *testing.T) {
	l := DefaultLogger
	l = WithReplace(l)
	l = WithReplace(l)
	l = WithReplace(l, "unpaired_key")
	c, ok := l.(*logger)
	if !ok || len(c.prefix) != 2 {
		t.Error("test failed")
	}

	l = WithReplace(l, "ts", DefaultTimestamp, "caller", DefaultCaller, "unpaired_key")
	l = WithReplace(l, "ts", Timestamp(time.ANSIC), "caller", Caller(-1), "unpaired_key")
	c, ok = l.(*logger)
	if !ok || len(c.prefix) != 6 {
		t.Error("test failed")
	}

	_ = l.Log(LevelInfo, "key1", "value1")
}
