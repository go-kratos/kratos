package log

import (
	"fmt"
	"os"
	"sync"
)

// globalLogger is designed as a global logger in current process.
var (
	global = &loggerAppliance{}
)

type loggerAppliance struct {
	lock   sync.Mutex
	logger Logger
	helper *Helper
}

func init() {
	global.SetLogger(DefaultLogger)
}

func (a *loggerAppliance) SetLogger(in Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.logger = in
	a.helper = NewHelper(a.logger)
}

func (a *loggerAppliance) GetLogger() Logger {
	return a.logger
}

// SetLogger should be called before any other log call.
// And it is NOT THREAD SAFE.
func SetLogger(logger Logger) {
	global.SetLogger(logger)
}

func GetLogger() Logger {
	return global.GetLogger()
}

// Log Print log by level and keyvals.
func Log(level Level, keyvals ...interface{}) {
	_ = global.logger.Log(level, keyvals...)
}

// Debug logs a message at debug level.
func Debug(a ...interface{}) {
	global.helper.Log(LevelDebug, global.helper.msgKey, fmt.Sprint(a...))
}

// Debugf logs a message at debug level.
func Debugf(format string, a ...interface{}) {
	global.helper.Log(LevelDebug, global.helper.msgKey, fmt.Sprintf(format, a...))
}

// Debugw logs a message at debug level.
func Debugw(keyvals ...interface{}) {
	global.helper.Log(LevelDebug, keyvals...)
}

// Info logs a message at info level.
func Info(a ...interface{}) {
	global.helper.Log(LevelInfo, global.helper.msgKey, fmt.Sprint(a...))
}

// Infof logs a message at info level.
func Infof(format string, a ...interface{}) {
	global.helper.Log(LevelInfo, global.helper.msgKey, fmt.Sprintf(format, a...))
}

// Infow logs a message at info level.
func Infow(keyvals ...interface{}) {
	global.helper.Log(LevelInfo, keyvals...)
}

// Warn logs a message at warn level.
func Warn(a ...interface{}) {
	global.helper.Log(LevelWarn, global.helper.msgKey, fmt.Sprint(a...))
}

// Warnf logs a message at warnf level.
func Warnf(format string, a ...interface{}) {
	global.helper.Log(LevelWarn, global.helper.msgKey, fmt.Sprintf(format, a...))
}

// Warnw logs a message at warnf level.
func Warnw(keyvals ...interface{}) {
	global.helper.Log(LevelWarn, keyvals...)
}

// Error logs a message at error level.
func Error(a ...interface{}) {
	global.helper.Log(LevelError, global.helper.msgKey, fmt.Sprint(a...))
}

// Errorf logs a message at error level.
func Errorf(format string, a ...interface{}) {
	global.helper.Log(LevelError, global.helper.msgKey, fmt.Sprintf(format, a...))
}

// Errorw logs a message at error level.
func Errorw(keyvals ...interface{}) {
	global.helper.Log(LevelError, keyvals...)
}

// Fatal logs a message at fatal level.
func Fatal(a ...interface{}) {
	global.helper.Log(LevelFatal, global.helper.msgKey, fmt.Sprint(a...))
	os.Exit(1)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, a ...interface{}) {
	global.helper.Log(LevelFatal, global.helper.msgKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// Fatalw logs a message at fatal level.
func Fatalw(keyvals ...interface{}) {
	global.helper.Log(LevelFatal, keyvals...)
	os.Exit(1)
}
