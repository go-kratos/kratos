package log

// Logger is a logger interface.
type Logger interface {
	Print(kvpair ...interface{})
}

type prefix struct {
	log    Logger
	kvpair []interface{}
}

func (l *prefix) Print(kvpair ...interface{}) {
	l.log.Print(append(l.kvpair, kvpair...)...)
}

// WithPrefix .
func WithPrefix(l Logger, kvpair ...interface{}) Logger {
	return &prefix{log: l, kvpair: kvpair}
}

// Debug .
func Debug(l Logger) Logger {
	return WithPrefix(l, LevelKey, LevelDebug)
}

// Info .
func Info(l Logger) Logger {
	return WithPrefix(l, LevelKey, LevelInfo)
}

// Warn .
func Warn(l Logger) Logger {
	return WithPrefix(l, LevelKey, LevelWarn)
}

// Error .
func Error(l Logger) Logger {
	return WithPrefix(l, LevelKey, LevelError)
}
