package slog

import (
	"fmt"
	"log/slog"

	"github.com/go-kratos/kratos/v2/log"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	log *slog.Logger
}

func NewLogger(slog *slog.Logger) log.Logger {
	return &Logger{slog}
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	keylen := len(keyvals)
	if keylen == 0 || keylen%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug("", keyvals...)
	case log.LevelInfo:
		l.log.Info("", keyvals...)
	case log.LevelWarn:
		l.log.Warn("", keyvals...)
	case log.LevelError:
		l.log.Error("", keyvals...)
	case log.LevelFatal:
		l.log.Error("", keyvals...)
	}
	return nil
}
