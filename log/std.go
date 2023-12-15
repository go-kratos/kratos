package log

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

var _ Logger = (*stdLogger)(nil)

// stdLogger corresponds to the standard library's [log.Logger] and provides
// similar capabilities. It also can be used concurrently by multiple goroutines.
type stdLogger struct {
	w         io.Writer
	isDiscard bool
	mu        sync.Mutex
	pool      *sync.Pool
}

// NewStdLogger new a logger with writer.
func NewStdLogger(w io.Writer) Logger {
	return &stdLogger{
		w:         w,
		isDiscard: w == io.Discard,
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

// Log print the kv pairs log.
func (l *stdLogger) Log(level Level, keyvals ...interface{}) error {
	if l.isDiscard || len(keyvals) == 0 {
		return nil
	}
	if (len(keyvals) & 1) == 1 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}

	buf := l.pool.Get().(*bytes.Buffer)
	defer l.pool.Put(buf)

	buf.WriteString(level.String())
	for i := 0; i < len(keyvals); i += 2 {
		_, _ = fmt.Fprintf(buf, " %s=%v", keyvals[i], keyvals[i+1])
	}
	buf.WriteByte('\n')
	defer buf.Reset()

	l.mu.Lock()
	defer l.mu.Unlock()
	_, err := l.w.Write(buf.Bytes())
	return err
}

func (l *stdLogger) Close() error {
	return nil
}
