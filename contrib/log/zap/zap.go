package zap

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/go-kratos/kratos/v2/log"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	log    *zap.Logger
	msgKey string
}

type Option func(*Logger)

// WithMessageKey with message key.
func WithMessageKey(key string) Option {
	return func(l *Logger) {
		l.msgKey = key
	}
}

func NewLogger(zlog *zap.Logger) *Logger {
	return &Logger{
		log:    zlog,
		msgKey: log.DefaultMessageKey,
	}
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	// If logging at this level is completely disabled, skip the overhead of
	// string formatting.
	if zapcore.Level(level) < zapcore.DPanicLevel && !l.log.Core().Enabled(zapcore.Level(level)) {
		return nil
	}
	var (
		msg    = ""
		keylen = len(keyvals)
	)
	if keylen == 0 || keylen%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	data := make([]zap.Field, 0, (keylen/2)+1)
	for i := 0; i < keylen; i += 2 {
		if keyvals[i].(string) == l.msgKey {
			msg, _ = keyvals[i+1].(string)
			continue
		}
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug(msg, data...)
	case log.LevelInfo:
		l.log.Info(msg, data...)
	case log.LevelWarn:
		l.log.Warn(msg, data...)
	case log.LevelError:
		l.log.Error(msg, data...)
	case log.LevelFatal:
		l.log.Fatal(msg, data...)
	}
	return nil
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

func (l *Logger) Close() error {
	return l.Sync()
}
