package log

import (
	"runtime"
	"strconv"
	"strings"
)

// Valuer is returns a log value.
type Valuer func() interface{}

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
