package log

type nopLogger struct{}

func (n *nopLogger) Print(kvpair ...interface{}) {}

func (n *nopLogger) Close() error { return nil }
