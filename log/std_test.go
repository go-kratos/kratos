package log

import "testing"

func TestStdLogger(t *testing.T) {
	logger := DefaultLogger
	logger = With(logger, "caller", DefaultCaller, "ts", DefaultTimestamp)

	_ = logger.Log(LevelInfo, "msg", "test debug")
	_ = logger.Log(LevelInfo, "msg", "test info")
	_ = logger.Log(LevelInfo, "msg", "test warn")
	_ = logger.Log(LevelInfo, "msg", "test error")
	_ = logger.Log(LevelDebug, "singular")

	logger2 := DefaultLogger
	_ = logger2.Log(LevelDebug)
}
