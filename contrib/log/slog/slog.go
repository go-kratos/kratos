package slog

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"slices"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	handler slog.Handler
}

func NewLogger(handler slog.Handler) log.Logger {
	return &Logger{handler}
}

func (l *Logger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if l == nil || !l.handler.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(3, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	_ = l.handler.Handle(ctx, r)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(context.Background(), slog.LevelDebug, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(context.Background(), slog.LevelWarn, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(context.Background(), slog.LevelInfo, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log(context.Background(), slog.LevelError, msg, args...)
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	keylen := len(keyvals)
	if keylen == 0 || keylen%2 != 0 {
		l.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	msg, keyvals := extractMessage(keyvals)

	switch level {
	case log.LevelDebug:
		l.Debug(msg, keyvals...)
	case log.LevelInfo:
		l.Info(msg, keyvals...)
	case log.LevelWarn:
		l.Warn(msg, keyvals...)
	case log.LevelError:
		l.Error(msg, keyvals...)
	case log.LevelFatal:
		l.Error(msg, keyvals...)
	}
	return nil
}

func extractMessage(keyvals []interface{}) (msg string, kvs []interface{}) {
	for i := 1; i < len(keyvals); i++ {
		if keyvals[i-1] == log.DefaultMessageKey {
			return fmt.Sprint(keyvals[i]), slices.Delete(keyvals, i-1, i+1)
		}
	}
	return "", keyvals
}
