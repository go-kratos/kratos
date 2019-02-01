package log

import (
	"context"
	"fmt"
)

// Config log config.
type Config struct {
	Family string
	Host   string

	// stdout
	Stdout bool

	// file
	Dir string
	// buffer size
	FileBufferSize int64
	// MaxLogFile
	MaxLogFile int
	// RotateSize
	RotateSize int64

	// V Enable V-leveled logging at the specified level.
	V int32
	// Module=""
	// The syntax of the argument is a map of pattern=N,
	// where pattern is a literal file name (minus the ".go" suffix) or
	// "glob" pattern and N is a V level. For instance:
	// [module]
	//   "service" = 1
	//   "dao*" = 2
	// sets the V level to 2 in all Go files whose names begin "dao".
	Module map[string]int32
	// Filter tell log handler which field are sensitive message, use * instead.
	Filter []string
}

var (
	h Handler
	c *Config
)

// D represents a map of entry level data used for structured logging.
// type D map[string]interface{}
type D struct {
	Key   string
	Value interface{}
}

// KV return a log kv for logging field.
func KV(key string, value interface{}) D {
	return D{
		Key:   key,
		Value: value,
	}
}

// Info logs a message at the info log level.
func Info(format string, args ...interface{}) {
	h.Log(context.Background(), _infoLevel, KV(_log, fmt.Sprintf(format, args...)))
}

// Warn logs a message at the warning log level.
func Warn(format string, args ...interface{}) {
	h.Log(context.Background(), _warnLevel, KV(_log, fmt.Sprintf(format, args...)))
}

// Error logs a message at the error log level.
func Error(format string, args ...interface{}) {
	h.Log(context.Background(), _errorLevel, KV(_log, fmt.Sprintf(format, args...)))
}

func logw(args []interface{}) []D {
	if len(args)%2 != 0 {
		Warn("log: the variadic must be plural, the last one will ignored")
	}
	ds := make([]D, 0, len(args)/2)
	for i := 0; i < len(args)-1; i = i + 2 {
		if key, ok := args[i].(string); ok {
			ds = append(ds, KV(key, args[i+1]))
		} else {
			Warn("log: key must be string, get %T, ignored", args[i])
		}
	}
	return ds
}
