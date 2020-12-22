package stdlog

import (
	"bytes"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

var _ log.Logger = (*stdLogger)(nil)

// Option is std logger option.
type Option func(*options)

type options struct {
	prefix string
	flag   int
	skip   int
	out    io.Writer
}

// Prefix with logger prefix.
func Prefix(prefix string) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}

// Flag with logger flag.
func Flag(flag int) Option {
	return func(o *options) {
		o.flag = flag
	}
}

// Skip with caller skip.
func Skip(skip int) Option {
	return func(o *options) {
		o.skip = skip
	}
}

// Writer with logger writer.
func Writer(out io.Writer) Option {
	return func(o *options) {
		o.out = out
	}
}

type stdLogger struct {
	opts options
	log  *stdlog.Logger
	pool *sync.Pool
}

// NewLogger new a std logger with options.
func NewLogger(opts ...Option) log.Logger {
	options := options{
		flag: stdlog.LstdFlags,
		skip: 2,
		out:  os.Stdout,
	}
	for _, o := range opts {
		o(&options)
	}
	return &stdLogger{
		opts: options,
		log:  stdlog.New(options.out, options.prefix, options.flag),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func stackTrace(path string) string {
	idx := strings.LastIndexByte(path, '/')
	if idx == -1 {
		return path
	}
	idx = strings.LastIndexByte(path[:idx], '/')
	if idx == -1 {
		return path
	}
	return path[idx+1:]
}

func (s *stdLogger) Print(kvpair ...interface{}) {
	if len(kvpair) == 0 {
		return
	}
	if len(kvpair)%2 != 0 {
		kvpair = append(kvpair, "")
	}
	buf := s.pool.Get().(*bytes.Buffer)
	if _, file, line, ok := runtime.Caller(s.opts.skip); ok {
		buf.WriteString(fmt.Sprintf("source=%s:%d ", stackTrace(file), line))
	}
	for i := 0; i < len(kvpair); i += 2 {
		fmt.Fprintf(buf, "%s=%s ", kvpair[i], kvpair[i+1])
	}
	s.log.Println(buf.String())
	buf.Reset()
	s.pool.Put(buf)
}
