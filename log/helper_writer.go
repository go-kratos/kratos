package log

import "io"

type writerWrapper struct {
	helper *Helper
	level  Level
}

type WriterOptionFn func(w *writerWrapper)

// WithWriteLevel set writerWrapper level.
func WithWriterLevel(level Level) WriterOptionFn {
	return func(w *writerWrapper) {
		w.level = level
	}
}

// WithWriteMessageKey set writerWrapper helper message key.
func WithWriteMessageKey(key string) WriterOptionFn {
	return func(w *writerWrapper) {
		w.helper.msgKey = key
	}
}

// NewWriter return a writer wrapper.
func NewWriter(logger Logger, opts ...WriterOptionFn) io.Writer {
	ww := &writerWrapper{
		helper: NewHelper(logger, WithMessageKey(DefaultMessageKey)),
		level:  LevelInfo, // default level
	}
	for _, opt := range opts {
		opt(ww)
	}
	return ww
}

func (ww *writerWrapper) Write(p []byte) (int, error) {
	ww.helper.Log(ww.level, ww.helper.msgKey, string(p))
	return 0, nil
}
