package log

import (
	"context"
	"testing"
	"time"
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

func TestWith(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp, "caller", DefaultCaller, "test_error_key_val_pair")
	logger = With(logger, "ts", Timestamp(time.ANSIC), "caller", Caller(-1), "test_error_key_val_pair")
	_ = logger.Log(LevelInfo, "key1", "value1")
}
