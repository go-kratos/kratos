package log

import (
	"context"
	"reflect"
	"testing"
)

func TestUniqKeys(t *testing.T) {
	uniq := uniqKeys([]interface{}{"k1", "v1", "k1", "vv1", "k1", "vvv1", "k2", "v2"})
	if !reflect.DeepEqual(uniq, []interface{}{"k1", "vvv1", "k2", "v2"}) {
		t.Error("inconsistent data after deduplication")
	}
}

func TestInfo(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp)
	logger = With(logger, "caller", DefaultCaller)
	logger = With(logger, "caller", Caller(1))
	_ = logger.Log(LevelInfo, "key1", "value1")
}

func TestWithContext(t *testing.T) {
	WithContext(context.Background(), nil)
}
