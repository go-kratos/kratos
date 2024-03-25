package log

import (
	"testing"
)

func TestInfo(_ *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp)
	logger = With(logger, "caller", DefaultCaller)
	_ = logger.Log(LevelInfo, "key1", "value1")
}

func TestLoggerFunc(t *testing.T) {
	type item struct {
		level Level
		kvs   []interface{}
	}

	var ch = make(chan item, 1) //nolint:gofumpt
	logger := LoggerFunc(func(level Level, keyvals ...interface{}) error {
		ch <- item{level: level, kvs: keyvals}
		return nil
	})
	if logger.Log(LevelInfo, "key1", "value1") != nil {
		t.Fatal("expect nil")
	}

	i := <-ch
	if i.level != LevelInfo {
		t.Fatal("expect LevelInfo")
	}
	if len(i.kvs) != 2 {
		t.Fatal("expect 2")
	}
	if i.kvs[0] != "key1" || i.kvs[1] != "value1" {
		t.Fatal("expect key1=value1")
	}
}
