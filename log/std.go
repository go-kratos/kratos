package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sync"
)

var _ Logger = (*stdLogger)(nil)

type stdLogger struct {
	log  *log.Logger
	pool *sync.Pool
}

// NewStdLogger new a std logger with options.
func NewStdLogger(w io.Writer) Logger {
	return &stdLogger{
		log: log.New(w, "", 0),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Log print the kv pairs log.
func (l *stdLogger) Log(kv ...interface{}) error {
	if len(kv) == 0 {
		return nil
	}
	if len(kv)%2 != 0 {
		kv = append(kv, "")
	}
	buf := l.pool.Get().(*bytes.Buffer)
	for i := 0; i < len(kv); i += 2 {
		fmt.Fprintf(buf, "%s=%v ", kv[i], kv[i+1])
	}
	l.log.Output(4, buf.String())
	buf.Reset()
	l.pool.Put(buf)
	return nil
}
