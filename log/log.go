package log

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	defaultDepth = 4
	// DefaultCaller is a Valuer that returns the file, line and func name.
	DefaultCaller = Caller(defaultDepth)
	// DefaultCallerFuncName is a Valuer that returns the func name.
	DefaultCallerFuncName = CallerFuncName(defaultDepth)
	// DefaultCallerFile is a Valuer that returns the file.
	DefaultCallerFile = CallerFile(defaultDepth)
	// DefaultCallerLine is a Valuer that returns the line.
	DefaultCallerLine = CallerLine(defaultDepth)

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

// Caller returns a Valuer that returns a pkg/file:line,name description of the caller.
func Caller(skip int) Valuer {
	return func(context.Context) interface{} {
		pc, file, line, _ := callerDepth(skip)
		idx := strings.LastIndexByte(file, '/')
		name := runtime.FuncForPC(pc).Name()
		index := strings.LastIndex(name, "/")
		return file[idx+1:] + ":" + strconv.Itoa(line) + "," + name[index+1:]
	}
}

// CallerFuncName returns a Valuer that returns a package.func description of the caller.
func CallerFuncName(skip int) Valuer {
	return func(ctx context.Context) interface{} {
		fn, _, _, ok := callerDepth(skip)
		if ok {
			name := runtime.FuncForPC(fn).Name()
			index := strings.LastIndexByte(name, '/')
			return name[index+1:]
		}
		return ""
	}
}

// CallerFile returns a Valuer that returns a pkg/file description of the caller.
func CallerFile(skip int) Valuer {
	return func(ctx context.Context) interface{} {
		_, file, _, ok := callerDepth(skip)
		if ok {
			return file[strings.LastIndex(file, "/")+1:]
		}
		return ""
	}
}

// CallerLine returns a Valuer that returns a line description of the caller.
func CallerLine(skip int) Valuer {
	return func(ctx context.Context) interface{} {
		_, _, line, _ := callerDepth(skip)
		return line
	}
}

// CallerDepth returns skip caller.
func callerDepth(skip int) (uintptr, string, int, bool) {
	d := skip
	pc, file, line, ok := runtime.Caller(skip)
	if strings.LastIndex(file, "/log/filter.go") > 0 {
		d++
		pc, file, line, ok = runtime.Caller(d)
	}
	if strings.LastIndex(file, "/log/helper.go") > 0 {
		d++
		pc, file, line, ok = runtime.Caller(d)
	}
	return pc, file, line, ok
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
