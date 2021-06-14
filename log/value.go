package log

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"
)

var (
	// DefaultCaller is a Valuer that returns the file and line.
	DefaultCaller = Caller(3)

	// DefaultTimestamp is a Valuer that returns the current wallclock time.
	DefaultTimestamp = Timestamp(time.RFC3339)
)

// Valuer is returns a log value.
type Valuer func(ctx context.Context) interface{}

// Value return the function value.
func Value(ctx context.Context, v interface{}) interface{} {
	if v, ok := v.(Valuer); ok {
		return v(ctx)
	}
	return v
}

// Caller returns returns a Valuer that returns a pkg/file:line description of the caller.
func Caller(depth int) Valuer {
	return func(context.Context) interface{} {
		_, file, line, _ := runtime.Caller(depth)
		if strings.LastIndex(file, "/log/filter.go") > 0 {
			depth++
			_, file, line, _ = runtime.Caller(depth)
		}
		if strings.LastIndex(file, "/log/helper.go") > 0 {
			depth++
			_, file, line, _ = runtime.Caller(depth)
		}
		idx := strings.LastIndexByte(file, '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}

// Timestamp returns a timestamp Valuer with a custom time format.
func Timestamp(layout string) Valuer {
	return func(context.Context) interface{} {
		return time.Now().Format(layout)
	}
}

// TraceID returns a traceid valuer.
func TraceID() Valuer {
	return func(ctx context.Context) interface{} {
		if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
			return span.TraceID().String()
		}
		return ""
	}
}

// SpanID returns a spanid valuer.
func SpanID() Valuer {
	return func(ctx context.Context) interface{} {
		if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
			return span.SpanID().String()
		}
		return ""
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
