package log

type nopLogger struct{}

func (n *nopLogger) Print(level Level, kvpair ...interface{}) {}
