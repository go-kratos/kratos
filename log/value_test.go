package log

import (
	"context"
	"reflect"
	"testing"
)

func TestValue(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	_ = logger.Log(LevelInfo, "msg", "helloworld")

	logger2 := DefaultLogger
	logger2 = With(logger2)
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
	if reflect.DeepEqual(got, 3) {
		t.Errorf("Value() = %v, want %v", got, 3)
	}
}
