package log

import (
	"testing"
)

func TestLogger(t *testing.T) {
	logger := DefaultLogger
	logger.Print(LevelInfo, "log", "test debug")
}
