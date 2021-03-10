package log

import "os"

var (
	// DefaultLogger is default logger.
	DefaultLogger Logger = NewStdLogger(os.Stderr)
)

// Logger is a logger interface.
type Logger interface {
	Print(level Level, pairs ...interface{})
}

type logger struct {
	log   Logger
	pairs []interface{}
}

func (l *logger) Print(level Level, pairs ...interface{}) {
	l.log.Print(level, append(pairs, l.pairs...)...)
}

// With with logger kv pairs.
func With(log Logger, pairs ...interface{}) Logger {
	if len(pairs) == 0 {
		return log
	}
	return &logger{log: log, pairs: pairs}
}
