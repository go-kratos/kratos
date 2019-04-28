package grocery

type level int

const (
	FINEST level = iota
	FINE
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
)

var (
	kLevelStrings = [...]string{"FINEST", "FINE", "DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL"}
)

func (l level) String() string {
	if l < 0 || int(l) > len(kLevelStrings) {
		return "UNKNOWN"
	}
	return kLevelStrings[int(l)]
}

type LogRecord struct {
	Level   level  // The log level
	Message string // The log message
}

func (r *LogRecord) String() string {
	return r.Message
}
