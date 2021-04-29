package log

import (
	"runtime"
	"strconv"
	"strings"
)

// Valuer is returns a log value.
type Valuer func() interface{}

func bindValues(keyvals []interface{}) {
	for i := 1; i < len(keyvals); i += 2 {
		if v, ok := keyvals[i].(Valuer); ok {
			keyvals[i] = v()
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

// Value return the function value.
func Value(v interface{}) interface{} {
	if v, ok := v.(Valuer); ok {
		return v()
	}
	return v
}

// Caller returns returns a Valuer that returns a pkg/file:line description of the caller.
func Caller(depth int) Valuer {
	return func() interface{} {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			return nil
		}
		idx := strings.LastIndexByte(file, '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}
