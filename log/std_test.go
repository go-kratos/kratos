package log

import "testing"

func TestFmtLogger(t *testing.T) {
	logger := DefaultLogger
	logger.Print(LevelInfo, "log", "test debug")
}
