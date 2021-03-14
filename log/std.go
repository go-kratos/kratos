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
		log: log.New(w, "", log.LstdFlags),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Print print the kv pairs log.
func (l *stdLogger) Print(pairs ...interface{}) {
	if len(pairs) == 0 {
		return
	}
	if len(pairs)%2 != 0 {
		pairs = append(pairs, "")
	}
	buf := l.pool.Get().(*bytes.Buffer)
	for i := 0; i < len(pairs); i += 2 {
		fmt.Fprintf(buf, "%s=%v ", pairs[i], Value(pairs[i+1]))
	}
	l.log.Output(4, buf.String())
	buf.Reset()
	l.pool.Put(buf)
}
