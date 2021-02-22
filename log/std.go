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
func (s *stdLogger) Print(kvpair ...interface{}) {
	if len(kvpair) == 0 {
		return
	}
	if len(kvpair)%2 != 0 {
		kvpair = append(kvpair, "")
	}
	buf := s.pool.Get().(*bytes.Buffer)
	for i := 0; i < len(kvpair); i += 2 {
		fmt.Fprintf(buf, "%s=%v ", kvpair[i], kvpair[i+1])
	}
	s.log.Println(buf.String())
	buf.Reset()
	s.pool.Put(buf)
}
