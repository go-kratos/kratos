package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
)

type stdlog struct {
	log *log.Logger
}

// NewStdLogger .
func NewStdLogger(out io.Writer) Logger {
	return &stdlog{log: log.New(out, "", log.LstdFlags)}
}

func (s *stdlog) Print(kvpair ...interface{}) {
	if len(kvpair) == 0 {
		return
	}
	if len(kvpair)%2 != 0 {
		kvpair = append(kvpair, "")
	}
	buf := &bytes.Buffer{}
	for i := 0; i < len(kvpair); i += 2 {
		fmt.Fprintf(buf, "%s=%s ", kvpair[i], kvpair[i+1])
	}
	s.log.Println(buf.String())
}
