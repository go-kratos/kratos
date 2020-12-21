package log

// level is a logger level.
type level int8

const (
	// LevelDebug is logger debug level.
	LevelDebug level = iota
	// LevelInfo is logger info level.
	LevelInfo
	// LevelWarn is logger warn level.
	LevelWarn
	// LevelError is logger error level.
	LevelError
)

// LevelKey is logger level key.
const LevelKey = "level"

// Enabled .
func (l level) Enabled(lv level) bool {
	return lv >= l
}

func (l level) String() string {
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
