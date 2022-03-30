package log

import (
	"bytes"
	"io"
	"testing"
)

func TestFilterAll(t *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	log := NewHelper(NewFilter(logger,
		FilterLevel(LevelDebug),
		FilterKey("username"),
		FilterValue("hello"),
		FilterFunc(testFilterFunc),
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
	log := NewFilter(logger)
	_ = log.Log(LevelDebug, "msg1", "te1st debug")
	logHelper := NewHelper(NewFilter(logger))
	logHelper.Log(LevelDebug, "msg1", "te1st debug")
}

func TestFilterKey(t *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	log := NewHelper(NewFilter(logger, FilterKey("password")))
	log.Debugw("password", "123456")
}

func TestFilterValue(t *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	log := NewHelper(NewFilter(logger, FilterValue("debug")))
	log.Debugf("test %s", "debug")
}

func TestFilterFunc(t *testing.T) {
	logger := With(DefaultLogger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	log := NewHelper(NewFilter(logger, FilterFunc(testFilterFunc)))
	log.Debug("debug level")
	log.Infow("password", "123456")
}

func BenchmarkFilterKey(b *testing.B) {
	log := NewHelper(NewFilter(NewStdLogger(io.Discard), FilterKey("password")))
	for i := 0; i < b.N; i++ {
		log.Infow("password", "123456")
	}
}

func BenchmarkFilterValue(b *testing.B) {
	log := NewHelper(NewFilter(NewStdLogger(io.Discard), FilterValue("password")))
	for i := 0; i < b.N; i++ {
		log.Infow("password")
	}
}

func BenchmarkFilterFunc(b *testing.B) {
	log := NewHelper(NewFilter(NewStdLogger(io.Discard), FilterFunc(testFilterFunc)))
	for i := 0; i < b.N; i++ {
		log.Info("password", "123456")
	}
}

func testFilterFunc(level Level, keyvals ...interface{}) bool {
	if level == LevelWarn {
		return true
	}
	for i := 0; i < len(keyvals); i++ {
		if keyvals[i] == "password" {
			keyvals[i+1] = fuzzyStr
		}
	}
	return false
}

func TestFilterFuncWitchLoggerPrefix(t *testing.T) {
	buf := new(bytes.Buffer)
	tests := []struct {
		logger Logger
		want   string
	}{
		{
			logger: NewFilter(With(NewStdLogger(buf), "caller", "caller", "prefix", "whaterver"), FilterFunc(testFilterFuncWithLoggerPrefix)),
			want:   "",
		},
		{
			logger: NewFilter(With(NewStdLogger(buf), "caller", "caller"), FilterFunc(testFilterFuncWithLoggerPrefix)),
			want:   "INFO caller=caller msg=msg\n",
		},
	}

	for _, tt := range tests {
		err := tt.logger.Log(LevelInfo, "msg", "msg")
		if err != nil {
			t.Fatal("err should be nil")
		}
		got := buf.String()
		if got != tt.want {
			t.Fatalf("filter should catch prefix, want %s, got %s.", tt.want, got)
		}
		buf.Reset()
	}
}

func testFilterFuncWithLoggerPrefix(level Level, keyvals ...interface{}) bool {
	if level == LevelWarn {
		return true
	}
	for i := 0; i < len(keyvals); i += 2 {
		if keyvals[i] == "prefix" {
			return true
		}
	}
	return false
}
