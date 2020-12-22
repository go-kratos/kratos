package log

type nopLogger struct{}

func (n *nopLogger) Print(kvpair ...interface{}) {}
