package log

// Logger is a logger interface.
type Logger interface {
	Print(level Level, kvpair ...interface{})
	Close() error
}

type prefix struct {
	log    Logger
	kvpair []interface{}
}

func (l *prefix) Print(level Level, kvpair ...interface{}) {
	l.log.Print(level, []interface{}{l.kvpair, kvpair}...)
}

func (l *prefix) Close() error {
	return l.log.Close()
}

// With with logger kv pairs.
func With(l Logger, kvpair ...interface{}) Logger {
	return &prefix{log: l, kvpair: kvpair}
}
