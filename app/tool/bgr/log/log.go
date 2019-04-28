package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"go-common/app/tool/bgr/log/color"

	"golang.org/x/crypto/ssh/terminal"
)

// FDWriter interface extends io.Writer with file descriptor function
type FDWriter interface {
	io.Writer
	Fd() uintptr
}

// Logger struct definition
type Logger struct {
	mu    sync.RWMutex
	out   FDWriter
	color bool
	debug bool
	buf   strings.Builder
}

type prefix struct {
	Plain string
	Color string
}

const (
	_plainError = "[ ERROR ] "
	_plainWarn  = "[ WARN ]  "
	_plainInfo  = "[ INFO ]  "
	_plainDebug = "[ DEBUG ] "
	_plainFatal = "[ FATAL ] "
)

var (
	_prefixError = prefix{
		Plain: _plainError,
		Color: colorful.Red(_plainError),
	}

	_prefixWarn = prefix{
		Plain: _plainWarn,
		Color: colorful.Orange(_plainWarn),
	}

	_prefixInfo = prefix{
		Plain: _plainInfo,
		Color: colorful.Green(_plainInfo),
	}

	_prefixDebug = prefix{
		Plain: _plainDebug,
		Color: colorful.Purple(_plainDebug),
	}

	_prefixFatal = prefix{
		Plain: _plainFatal,
		Color: colorful.Gray(_plainFatal),
	}
)

// New returns new Logger instance with predefined writer output and
// automatically detect terminal coloring support
func New(out FDWriter, debug bool) *Logger {
	return &Logger{
		color: terminal.IsTerminal(int(out.Fd())),
		out:   out,
		debug: debug,
		buf:   strings.Builder{},
	}
}

func (l *Logger) output(prefix prefix, data string) (err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.buf.Reset()
	if l.color {
		if _, err = l.buf.WriteString(prefix.Color); err != nil {
			return
		}
	} else {
		if _, err = l.buf.WriteString(prefix.Plain); err != nil {
			return
		}
	}
	if _, err = l.buf.WriteString(data); err != nil {
		return
	}
	if data[len(data)-1] != '\n' {
		l.buf.WriteString("\n")
	}

	_, err = l.out.Write([]byte(l.buf.String()))
	return
}

// Error print error message to output
func (l *Logger) Error(v ...interface{}) {
	l.output(_prefixError, fmt.Sprintln(v...))
}

// Errorf print formatted error message to output
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.output(_prefixError, fmt.Sprintf(format, v...))
}

// Warn print warning message to output
func (l *Logger) Warn(v ...interface{}) {
	l.output(_prefixWarn, fmt.Sprintln(v...))
}

// Warnf print formatted warning message to output
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.output(_prefixWarn, fmt.Sprintf(format, v...))
}

// Info print informational message to output
func (l *Logger) Info(v ...interface{}) {
	l.output(_prefixInfo, fmt.Sprintln(v...))
}

// Infof print formatted informational message to output
func (l *Logger) Infof(format string, v ...interface{}) {
	l.output(_prefixInfo, fmt.Sprintf(format, v...))
}

// Debug print debug message to output if debug output enabled
func (l *Logger) Debug(v ...interface{}) {
	if l.debug {
		l.output(_prefixDebug, fmt.Sprintln(v...))
	}
}

// Debugf print formatted debug message to output if debug output enabled
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.debug {
		l.output(_prefixDebug, fmt.Sprintf(format, v...))
	}
}

// Fatal print fatal message to output and then exit(1)
func (l *Logger) Fatal(v ...interface{}) {
	l.output(_prefixFatal, fmt.Sprintln(v...))
	os.Exit(1)
}

// Fatalf print formatted fatal message to output and then exit(1)
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.output(_prefixFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}
