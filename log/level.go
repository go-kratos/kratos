package log

// Level is a logger level.
type Level int8

const (
	// LevelDebug is logger debug level.
	LevelDebug = iota
	// LevelInfo is logger info level.
	LevelInfo
	// LevelWarn is logger warn level.
	LevelWarn
	// LevelError is logger error level.
	LevelError
)

// Enabled .
func (l Level) Enabled(lv Level) bool {
	return lv >= l
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return ""
	}
}
