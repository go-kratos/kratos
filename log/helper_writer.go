package log

import "io"

type writerWrapper struct {
	helper *Helper
	level  Level
}

type writerOption func(w *writerWrapper)

// WithWriteLevel set writerWrapper level.
func WithWriterLevel(level Level) writerOption {
	return func(w *writerWrapper) {
		w.level = level
	}
}

// WithWriteMessageKey set writerWrapper helper message key.
func WithWriteMessageKey(key string) writerOption {
	return func(w *writerWrapper) {
		if key != "" {
			w.helper.msgKey = key
		}
	}
}

// NewWriter return a writer wrapper.
func NewWriter(logger Logger, opts ...writerOption) io.Writer {
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
