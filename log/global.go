package log

import (
	"sync"
)

// globalLogger is designed as a global logger in current process.
var global = &loggerAppliance{}

// loggerAppliance is the proxy of `Logger` to
// make logger change will affect all sub-logger.
type loggerAppliance struct {
	lock sync.Mutex
	Logger
	helper *Helper
}

func init() {
	global.SetLogger(DefaultLogger)
}

func (a *loggerAppliance) SetLogger(in Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Logger = in
	a.helper = NewHelper(a.Logger)
}

func (a *loggerAppliance) GetLogger() Logger {
	return a.Logger
}

// SetLogger should be called before any other log call.
// And it is NOT THREAD SAFE.
func SetLogger(logger Logger) {
	global.SetLogger(logger)
}

// GetLogger returns global logger appliance as logger in current process.
func GetLogger() Logger {
	return global
}

// Log Print log by level and keyvals.
func Log(level Level, keyvals ...interface{}) {
	global.helper.Log(level, keyvals...)
}

// Debug logs a message at debug level.
func Debug(a ...interface{}) {
	global.helper.Debug(a...)
}

// Debugf logs a message at debug level.
func Debugf(format string, a ...interface{}) {
	global.helper.Debugf(format, a...)
}

// Debugw logs a message at debug level.
func Debugw(keyvals ...interface{}) {
	global.helper.Debugw(keyvals...)
}

// Info logs a message at info level.
func Info(a ...interface{}) {
	global.helper.Info(a...)
}

// Infof logs a message at info level.
func Infof(format string, a ...interface{}) {
	global.helper.Infof(format, a...)
}

// Infow logs a message at info level.
func Infow(keyvals ...interface{}) {
	global.helper.Infow(keyvals...)
}

// Warn logs a message at warn level.
func Warn(a ...interface{}) {
	global.helper.Warn(a...)
}

// Warnf logs a message at warnf level.
func Warnf(format string, a ...interface{}) {
	global.helper.Warnf(format, a...)
}

// Warnw logs a message at warnf level.
func Warnw(keyvals ...interface{}) {
	global.helper.Warnw(keyvals...)
}

// Error logs a message at error level.
func Error(a ...interface{}) {
	global.helper.Error(a...)
}

// Errorf logs a message at error level.
func Errorf(format string, a ...interface{}) {
	global.helper.Errorf(format, a...)
}

// Errorw logs a message at error level.
func Errorw(keyvals ...interface{}) {
	global.helper.Errorw(keyvals...)
}

// Fatal logs a message at fatal level.
func Fatal(a ...interface{}) {
	global.helper.Fatal(a...)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, a ...interface{}) {
	global.helper.Fatalf(format, a...)
}

// Fatalw logs a message at fatal level.
func Fatalw(keyvals ...interface{}) {
	global.helper.Fatalw(keyvals...)
}
