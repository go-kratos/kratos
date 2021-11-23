package log

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	defaultDepth = 0
	// DefaultCaller is a Valuer that returns the external caller's file and line.
	DefaultCaller = Caller(defaultDepth)

	// baseDepth is the depth from logger.Log to Caller
	baseDepth = 2

	// DefaultTimestamp is a Valuer that returns the current wallclock time.
	DefaultTimestamp = Timestamp(time.RFC3339)
)

// relativeDepthKey is the key of depth from caller to logger.Log
type relativeDepthKey struct{}

// Valuer is returns a log value.
type Valuer func(ctx context.Context) interface{}

// Value return the function value.
func Value(ctx context.Context, v interface{}) interface{} {
	if v, ok := v.(Valuer); ok {
		return v(ctx)
	}
	return v
}

// Caller returns a Valuer that returns a pkg/file:line description of the caller.
func Caller(depth int) Valuer {
	return func(ctx context.Context) interface{} {
		relativeDepth := getRelativeDepth(ctx)
		_, file, line, _ := runtime.Caller(depth + relativeDepth + baseDepth)
		idx := strings.LastIndexByte(file, '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

// Set the relative depth of caller to logger.Log in ctx
func setRelativeDepth(ctx context.Context, relativeDepth int) context.Context {
	return context.WithValue(ctx, relativeDepthKey{}, relativeDepth)
}

// Get the relative depth of caller to logger.Log from ctx
func getRelativeDepth(ctx context.Context) int {
	if ctx != nil {
		if depth := ctx.Value(relativeDepthKey{}); depth != nil {
			return depth.(int)
		}
	}
	return 0
}

// Timestamp returns a timestamp Valuer with a custom time format.
func Timestamp(layout string) Valuer {
	return func(context.Context) interface{} {
		return time.Now().Format(layout)
	}
}

func bindValues(ctx context.Context, keyvals []interface{}) {
	for i := 1; i < len(keyvals); i += 2 {
		if v, ok := keyvals[i].(Valuer); ok {
			keyvals[i] = v(ctx)
		}
	}
}

func containsValuer(keyvals []interface{}) bool {
	for i := 1; i < len(keyvals); i += 2 {
		if _, ok := keyvals[i].(Valuer); ok {
			return true
		}
	}
	return false
}
