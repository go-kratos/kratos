package log

// Logger is a logger interface.
type Logger interface {
	Print(kvpair ...interface{})
}

type printer struct {
	log    Logger
	kvpair []interface{}
}

func newPrinter(log Logger, kvpair ...interface{}) *printer {
	return &printer{log: log, kvpair: kvpair}
}

func (l *printer) Print(kvpair ...interface{}) {
	l.log.Print(append(kvpair, l.kvpair...)...)
}

// With with logger kv pairs.
func With(log Logger, kvpair ...interface{}) Logger {
	return newPrinter(log, kvpair...)
}

// Debug returns a debug logger.
func Debug(log Logger) Logger {
	return newPrinter(log, LevelKey, LevelDebug)
}

// Info returns a info logger.
func Info(log Logger) Logger {
	return newPrinter(log, LevelKey, LevelInfo)
}

// Warn return a warn logger.
func Warn(log Logger) Logger {
	return newPrinter(log, LevelKey, LevelWarn)
}

// Error returns a error logger.
func Error(log Logger) Logger {
	return newPrinter(log, LevelKey, LevelError)
}
