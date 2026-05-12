package log

import (
	"log/slog"
	"strings"
)

// Level is a logger level.
type Level = slog.Level

// Leveler provides a log level.
type Leveler = slog.Leveler

// LevelVar is a variable log level.
type LevelVar = slog.LevelVar

// LevelKey is logger level key.
const LevelKey = slog.LevelKey

const (
	// LevelDebug is logger debug level.
	LevelDebug Level = slog.LevelDebug
	// LevelInfo is logger info level.
	LevelInfo Level = slog.LevelInfo
	// LevelWarn is logger warn level.
	LevelWarn Level = slog.LevelWarn
	// LevelError is logger error level.
	LevelError Level = slog.LevelError
	// LevelFatal is logger fatal level.
	LevelFatal Level = slog.LevelError + 4
)

// ParseLevel parses a level string into a logger Level value.
func ParseLevel(s string) Level {
	if strings.EqualFold(s, "FATAL") {
		return LevelFatal
	}
	var level slog.Level
	if err := level.UnmarshalText([]byte(s)); err == nil {
		return level
	}
	return LevelInfo
}
