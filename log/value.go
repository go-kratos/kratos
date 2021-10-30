package log

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	defaultDepth = 2
	// DefaultCaller is a Valuer that returns the file and line.
	DefaultCaller = Caller(0)

	// DefaultTimestamp is a Valuer that returns the current wallclock time.
	DefaultTimestamp = Timestamp(time.RFC3339)
)

type skipDepthKey struct{}

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
		curDepth := getSkipDepth(ctx)
		_, file, line, _ := runtime.Caller(depth + curDepth + defaultDepth)
		idx := strings.LastIndexByte(file, '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

// Set the skip depth for the ctx of the current logger
func setSkipDepth(ctx context.Context, depth int) context.Context {
	return context.WithValue(ctx, skipDepthKey{}, depth)
}

// Get the skipped depth from ctx
func getSkipDepth(ctx context.Context) int {
	if ctx != nil {
		if depth := ctx.Value(skipDepthKey{}); depth != nil {
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
