package log

type nopLogger struct{}

func (l *nopLogger) Print(kvpair ...interface{}) {}
