package log

import (
	"context"
	"testing"
)

func TestWithUniq(t *testing.T) {
	l := DefaultLogger
	l = WithUniq(l, "caller", 10)
	l = WithUniq(l, "caller", 11)
	c, ok := l.(*logger)
	if !ok {
		t.Error("Interface Logger is not implemented")
	}
	if v, ok := c.prefixUniq["caller"]; !ok {
		t.Error("prefix uniq map err")
	} else if c.prefix[v] != 11 {
		t.Error("prefix map does not save the last key")
	}
}

func TestInfo(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp)
	logger = WithUniq(logger, "caller", DefaultCaller)
	logger = WithUniq(logger, "caller", Caller(1))
	_ = logger.Log(LevelInfo, "key1", "value1")
}

func TestWithContext(t *testing.T) {
	WithContext(context.Background(), nil)
}
