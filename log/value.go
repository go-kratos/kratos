package log

import (
	"runtime"
	"strconv"
	"strings"
)

// Valuer is returns a log value.
type Valuer func() interface{}

func bindValues(pairs []interface{}) []interface{} {
	for i := 1; i < len(pairs); i += 2 {
		if v, ok := pairs[i].(Valuer); ok {
			pairs[i] = v()
		}
	}
	return pairs
}

// Caller returns a Valuer that returns a pkg/file:line description of the caller.
func Caller(skip int) Valuer {
	return func() interface{} {
		_, file, line, _ := runtime.Caller(skip)
		idx := strings.LastIndexByte(file, '/')
		return file[idx+1:] + ":" + strconv.Itoa(line)
	}
}
