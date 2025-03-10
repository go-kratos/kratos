package log

import (
	"context"
	"fmt"
	"os"
	"sync"
)

// globalLogger is designed as a global logger in current process.
var global = &loggerAppliance{}

// loggerAppliance is the proxy of `Logger` to
// make logger change will affect all sub-logger.
type loggerAppliance struct {
	lock sync.RWMutex
	Logger
}

func init() {
	global.SetLogger(DefaultLogger)
}

func (a *loggerAppliance) SetLogger(in Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Logger = in
}

// SetLogger should be called before any other log call.
// And it is NOT THREAD SAFE.
func SetLogger(logger Logger) {
	global.SetLogger(logger)
}

// GetLogger returns global logger appliance as logger in current process.
func GetLogger() Logger {
	global.lock.RLock()
	defer global.lock.RUnlock()
	return global.Logger
}

// Log Print log by level and keyvals.
func Log(level Level, keyvals ...any) {
	_ = global.Log(level, keyvals...)
}

// Context with context logger.
func Context(ctx context.Context) *Helper {
	return NewHelper(WithContext(ctx, global.Logger))
}

// Debug logs a message at debug level.
func Debug(a ...any) {
	_ = global.Log(LevelDebug, DefaultMessageKey, fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func Debugf(format string, a ...any) {
	_ = global.Log(LevelDebug, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func Debugw(keyvals ...any) {
	_ = global.Log(LevelDebug, keyvals...)
}

// Info logs a message at info level.
func Info(a ...any) {
	_ = global.Log(LevelInfo, DefaultMessageKey, fmt.Sprint(a...))
}

// Infof logs a message at info level.
func Infof(format string, a ...any) {
	_ = global.Log(LevelInfo, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func Infow(keyvals ...any) {
	_ = global.Log(LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func Warn(a ...any) {
	_ = global.Log(LevelWarn, DefaultMessageKey, fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func Warnf(format string, a ...any) {
	_ = global.Log(LevelWarn, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func Warnw(keyvals ...any) {
	_ = global.Log(LevelWarn, keyvals...)
}

// Error logs a message at error level.
func Error(a ...any) {
	_ = global.Log(LevelError, DefaultMessageKey, fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func Errorf(format string, a ...any) {
	_ = global.Log(LevelError, DefaultMessageKey, fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func Errorw(keyvals ...any) {
	_ = global.Log(LevelError, keyvals...)
}

// Fatal logs a message at fatal level.
func Fatal(a ...any) {
	_ = global.Log(LevelFatal, DefaultMessageKey, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, a ...any) {
	_ = global.Log(LevelFatal, DefaultMessageKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Fatalw logs a message at fatal level.
func Fatalw(keyvals ...any) {
	_ = global.Log(LevelFatal, keyvals...)
	os.Exit(1)
}
