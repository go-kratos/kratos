package log

import (
	"context"
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
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)

	filter := NewFilter(logger, FilterLevel(LevelDebug))

	helper := NewHelper(filter)

	logger2 := With(filter)

	logger3 := With(logger2)

	mLog := MultiLogger(logger, filter)

	logger4 := With(mLog)

	_ = logger.Log(LevelDebug, "msg", "value_test.go:47")
	_ = WithContext(context.Background(), logger).Log(LevelDebug, "msg", "value_test.go:48")
	_ = filter.Log(LevelDebug, "msg", "value_test.go:49")
	helper.Log(LevelDebug, "msg", "value_test.go:50")
	_ = logger2.Log(LevelDebug, "msg", "value_test.go:51")
	_ = logger3.Log(LevelDebug, "msg", "value_test.go:52")
	_ = mLog.Log(LevelDebug, "msg", "value_test.go:53")
	_ = logger4.Log(LevelDebug, "msg", "value_test.go:54")

	_ = mLog.Log(LevelDebug, "msg", "value_test.go:56")
	_ = logger3.Log(LevelDebug, "msg", "value_test.go:57")
	_ = logger2.Log(LevelDebug, "msg", "value_test.go:58")
	helper.Log(LevelDebug, "msg", "value_test.go:59")
	_ = filter.Log(LevelDebug, "msg", "value_test.go:60")
	_ = WithContext(context.Background(), logger).Log(LevelDebug, "msg", "value_test.go:61")
	_ = logger.Log(LevelDebug, "msg", "value_test.go:62")
}
