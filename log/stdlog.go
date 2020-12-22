package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sync"
)

type stdlog struct {
	log  *log.Logger
	pool *sync.Pool
}

// NewStdLogger .
func NewStdLogger(out io.Writer) Logger {
	return &stdlog{
		log: log.New(out, "", log.LstdFlags),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func (s *stdlog) Print(kvpair ...interface{}) {
	if len(kvpair) == 0 {
		return
	}
	if len(kvpair)%2 != 0 {
		kvpair = append(kvpair, "")
	}
	buf := s.pool.Get().(*bytes.Buffer)
	for i := 0; i < len(kvpair); i += 2 {
		fmt.Fprintf(buf, "%s=%s ", kvpair[i], kvpair[i+1])
	}
	s.log.Println(buf.String())
	buf.Reset()
	s.pool.Put(buf)
}
