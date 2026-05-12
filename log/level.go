package log

import "strings"

// Level is a logger level.
type Level int8

// LevelKey is logger level key.
const LevelKey = "level"

const (
	// LevelDebug is logger debug level.
	LevelDebug Level = iota - 1
	// LevelInfo is logger info level.
	LevelInfo
	// LevelWarn is logger warn level.
	LevelWarn
	// LevelError is logger error level.
	LevelError
	// LevelFatal is logger fatal level
	LevelFatal
)

const (
	levelDebugString = "DEBUG"
	levelInfoString  = "INFO"
	levelWarnString  = "WARN"
	levelErrorString = "ERROR"
	levelFatalString = "FATAL"
)

func (l Level) Key() string {
	return LevelKey
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return levelDebugString
	case LevelInfo:
		return levelInfoString
	case LevelWarn:
		return levelWarnString
	case LevelError:
		return levelErrorString
	case LevelFatal:
		return levelFatalString
	default:
		return ""
	}
}

// ParseLevel parses a level string into a logger Level value.
func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case levelDebugString:
		return LevelDebug
	case levelInfoString:
		return LevelInfo
	case levelWarnString:
		return LevelWarn
	case levelErrorString:
		return LevelError
	case levelFatalString:
		return LevelFatal
	}
	return LevelInfo
}
