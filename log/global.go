package log

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
)

// global is a place to store global logger.
var global atomic.Value

// loggerAppliance is a wrapper to make sure logger can be stored.
type loggerAppliance struct {
	logger Logger
}

func init() {
	SetLogger(DefaultLogger)
}

// SetLogger sets the global logger. This function should be called
// before any other log call.
func SetLogger(logger Logger) {
	global.Store(loggerAppliance{logger: logger})
}

// GetLogger returns global logger.
func GetLogger() Logger {
	return global.Load().(loggerAppliance).logger
}

// Log Print log by level and keyvals.
func Log(level Level, keyvals ...interface{}) {
	_ = GetLogger().Log(level, keyvals...)
}

// Context with context logger.
func Context(ctx context.Context) *Helper {
	return NewHelper(WithContext(ctx, GetLogger()))
}

// Debug logs a message at debug level.
func Debug(a ...interface{}) {
	_ = GetLogger().Log(LevelDebug, DefaultMessageKey, fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func Debugf(format string, a ...interface{}) {
	_ = GetLogger().Log(LevelDebug, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func Debugw(keyvals ...interface{}) {
	_ = GetLogger().Log(LevelDebug, keyvals...)
}

// Info logs a message at info level.
func Info(a ...interface{}) {
	_ = GetLogger().Log(LevelInfo, DefaultMessageKey, fmt.Sprint(a...))
}

// Infof logs a message at info level.
func Infof(format string, a ...interface{}) {
	_ = GetLogger().Log(LevelInfo, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func Infow(keyvals ...interface{}) {
	_ = GetLogger().Log(LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func Warn(a ...interface{}) {
	_ = GetLogger().Log(LevelWarn, DefaultMessageKey, fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func Warnf(format string, a ...interface{}) {
	_ = GetLogger().Log(LevelWarn, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func Warnw(keyvals ...interface{}) {
	_ = GetLogger().Log(LevelWarn, keyvals...)
}

// Error logs a message at error level.
func Error(a ...interface{}) {
	_ = GetLogger().Log(LevelError, DefaultMessageKey, fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func Errorf(format string, a ...interface{}) {
	_ = GetLogger().Log(LevelError, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func Errorw(keyvals ...interface{}) {
	_ = GetLogger().Log(LevelError, keyvals...)
}

// Fatal logs a message at fatal level.
func Fatal(a ...interface{}) {
	_ = GetLogger().Log(LevelFatal, DefaultMessageKey, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, a ...interface{}) {
	_ = GetLogger().Log(LevelFatal, DefaultMessageKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Fatalw logs a message at fatal level.
func Fatalw(keyvals ...interface{}) {
	_ = GetLogger().Log(LevelFatal, keyvals...)
	os.Exit(1)
}
