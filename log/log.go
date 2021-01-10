package log

import "sync"

var ppFree = sync.Pool{
	New: func() interface{} { return new(printer) },
}

// Logger is a logger interface.
type Logger interface {
	Print(kvpair ...interface{})
}

type printer struct {
	log     Logger
	kvpair  []interface{}
	recycle bool
}

func newPrinter() *printer {
	return ppFree.Get().(*printer)
}

func (l *printer) Print(kvpair ...interface{}) {
	l.log.Print(append(kvpair, l.kvpair...)...)
	if l.recycle {
		l.free()
	}
}

func (l *printer) free() {
	l.log = nil
	l.kvpair = nil
	ppFree.Put(l)
}

func with(l Logger, free bool, kvpair ...interface{}) Logger {
	p := newPrinter()
	p.log = l
	p.kvpair = kvpair
	p.recycle = free
	return p
}

// With with logger kv pairs.
func With(l Logger, kvpair ...interface{}) Logger {
	return with(l, false, kvpair)
}

// Debug returns a debug logger.
func Debug(l Logger) Logger {
	return with(l, true, LevelKey, LevelDebug)
}

// Info returns a info logger.
func Info(l Logger) Logger {
	return with(l, true, LevelKey, LevelInfo)
}

// Warn return a warn logger.
func Warn(l Logger) Logger {
	return with(l, true, LevelKey, LevelWarn)
}

// Error returns a error logger.
func Error(l Logger) Logger {
	return with(l, true, LevelKey, LevelError)
}
