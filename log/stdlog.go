package log

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
)

type stdLogger struct {
	log  *log.Logger
	pool *sync.Pool
}

func (l *stdLogger) Print(kvpair ...interface{}) {
	if len(kvpair) == 0 {
		return
	}
	if len(kvpair)%2 != 0 {
		kvpair = append(kvpair, "")
	}
	buf := l.pool.Get().(*bytes.Buffer)
	for i := 0; i < len(kvpair); i += 2 {
		fmt.Fprintf(buf, "%s=%s ", kvpair[i], kvpair[i+1])
	}
	l.log.Println(buf.String())
	buf.Reset()
	l.pool.Put(buf)
}

var defaultLogger Logger = &stdLogger{
	log: log.New(os.Stdout, "", log.LstdFlags),
	pool: &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	},
}

// SetLogger .
func SetLogger(logger Logger) {
	defaultLogger = logger
}

// GetLogger returns a logger instance with package name.
func GetLogger(module string) Logger {
	return WithPrefix(defaultLogger, "module", module)
}
