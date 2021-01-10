package log

// Logger is a logger interface.
type Logger interface {
	Print(level Level, kvpair ...interface{})
}

type suffix struct {
	log    Logger
	kvpair []interface{}
}

func (l *suffix) Print(level Level, kvpair ...interface{}) {
	l.log.Print(level, append(kvpair, l.kvpair...)...)
}

// With with logger kv pairs.
func With(l Logger, kvpair ...interface{}) Logger {
	return &suffix{log: l, kvpair: kvpair}
}
