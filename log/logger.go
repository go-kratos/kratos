package log

import (
	"context"
)

var defaultLogger = GetLogger("log")

// Logger is a logger interface.
type Logger interface {
	Verbose(v int) Verbose
	Print(ctx context.Context, lv Level, a ...interface{})
	Printf(ctx context.Context, lv Level, format string, a ...interface{})
	Printw(ctx context.Context, lv Level, kvpair ...interface{})
}

// SetLogger sets the logger that is used in application.
func SetLogger(logger Logger) {
	defaultLogger = logger
}

// Debug logs a message at debug level.
func Debug(ctx context.Context, a ...interface{}) {
	defaultLogger.Print(ctx, LevelDebug, a)
}

// Debugf logs a message at debug level.
func Debugf(ctx context.Context, format string, a ...interface{}) {
	defaultLogger.Printf(ctx, LevelDebug, format, a)
}

// Debugw logs a message at debug level.
func Debugw(ctx context.Context, kvpair ...interface{}) {
	defaultLogger.Printw(ctx, LevelDebug, kvpair)
}

// Info logs a message at info level.
func Info(ctx context.Context, a ...interface{}) {
	defaultLogger.Print(ctx, LevelInfo, a)
}

// Infof logs a message at info level.
func Infof(ctx context.Context, format string, a ...interface{}) {
	defaultLogger.Printf(ctx, LevelInfo, format, a)
}

// Infow logs a message at info level.
func Infow(ctx context.Context, kvpair ...interface{}) {
	defaultLogger.Printw(ctx, LevelInfo, kvpair)
}

// Warn logs a message at warn level.
func Warn(ctx context.Context, format string, a ...interface{}) {
	defaultLogger.Print(ctx, LevelWarn, a)
}

// Warnf logs a message at warnf level.
func Warnf(ctx context.Context, format string, a ...interface{}) {
	defaultLogger.Printf(ctx, LevelWarn, format, a)
}

// Warnw logs a message at warnf level.
func Warnw(ctx context.Context, kvpair ...interface{}) {
	defaultLogger.Printw(ctx, LevelWarn, kvpair)
}

// Error logs a message at error level.
func Error(ctx context.Context, a ...interface{}) {
	defaultLogger.Print(ctx, LevelWarn, a)
}

// Errorf logs a message at error level.
func Errorf(ctx context.Context, format string, a ...interface{}) {
	defaultLogger.Printf(ctx, LevelError, format, a)
}

// Errorw logs a message at error level.
func Errorw(ctx context.Context, kvpair ...interface{}) {
	defaultLogger.Printw(ctx, LevelError, kvpair)
}
