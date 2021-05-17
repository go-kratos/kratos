package log

import "testing"

func TestValue(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "ts", DefaultTimestamp, "caller", DefaultCaller)
	logger.Log(LevelInfo, "msg", "helloworld")
}
