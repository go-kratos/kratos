package log

import (
	"testing"
)

func TestFilterAll(t *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	log := NewHelper(NewFilter(logger,
		FilterLevel(LevelDebug),
		FilterKeys("username"),
		FilterValues("hello"),
		FilterHook(testFilterFunc),
	))
	log.Log(LevelDebug, "msg", "test debug")
	log.Info("hello")
	log.Infow("password", "123456")
	log.Infow("username", "kratos")
	log.Warn("warn log")
}
func TestFilterLevel(t *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	log := NewHelper(NewFilter(NewFilter(logger, FilterLevel(LevelWarn))))
	log.Log(LevelDebug, "msg1", "te1st debug")
	log.Debug("test debug")
	log.Debugf("test %s", "debug")
	log.Debugw("log", "test debug")
	log.Warn("warn log")
}

func TestFilterCaller(t *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	log := NewFilter(NewFilter(logger))
	log.Log(LevelDebug, "msg1", "te1st debug")
	logHelper := NewHelper(NewFilter(NewFilter(logger)))
	logHelper.Log(LevelDebug, "msg1", "te1st debug")
}

func TestFilterKey(t *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	log := NewHelper(NewFilter(logger, FilterKeys("password")))
	log.Debugw("password", "123456")
}

func TestFilterValue(t *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	log := NewHelper(NewFilter(logger, FilterValues("debug")))
	log.Debugf("test %s", "debug")
}

func TestFilterFunc(t *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	log := NewHelper(NewFilter(logger, FilterHook(testFilterFunc)))
	log.Debug("debug level")
	log.Infow("password", "123456")
}

func testFilterFunc(level Level, keyvals ...interface{}) bool {
	if level == LevelWarn {
		return true
	}
	for i := 0; i < len(keyvals); i++ {
		if keyvals[i] == "password" {
			keyvals[i+1] = "***"
		}
	}
	return false
}