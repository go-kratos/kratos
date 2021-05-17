package log

import "testing"

func TestStdLogger(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "caller", DefaultCaller, "ts", DefaultTimestamp)

	logger.Log(LevelInfo, "msg", "test debug")
	logger.Log(LevelInfo, "msg", "test info")
	logger.Log(LevelInfo, "msg", "test warn")
	logger.Log(LevelInfo, "msg", "test error")
}
