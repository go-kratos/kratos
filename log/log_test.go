package log

import (
	"context"
	"testing"
)

func TestInfo(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp)
	logger = With(logger, "caller", DefaultCaller)
	_ = logger.Log(LevelInfo, "key1", "value1")
}

func TestWithContext(t *testing.T) {
	WithContext(context.Background(), nil)
}
