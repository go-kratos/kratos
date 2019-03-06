package log

import (
	"context"
	"io"
	"os"
	"time"
)

const defaultPattern = "%L %d-%T %f %M"

var _defaultStdout = NewStdout()

// StdoutHandler stdout log handler
type StdoutHandler struct {
	out    io.Writer
	render Render
}

// NewStdout create a stdout log handler
func NewStdout() *StdoutHandler {
	return &StdoutHandler{
		out:    os.Stderr,
		render: newPatternRender(defaultPattern),
	}
}

// Log stdout loging, only for developing env.
func (h *StdoutHandler) Log(ctx context.Context, lv Level, args ...D) {
	d := make(map[string]interface{}, 10+len(args))
	for _, arg := range args {
		d[arg.Key] = arg.Value
	}
	// add extra fields
	addExtraField(ctx, d)
	d[_time] = time.Now().Format(_timeFormat)
	h.render.Render(h.out, d)
	h.out.Write([]byte("\n"))
}

// Close stdout loging
func (h *StdoutHandler) Close() error {
	return nil
}

// SetFormat set stdout log output format
// %T time format at "15:04:05.999"
// %t time format at "15:04:05"
// %D data format at "2006/01/02"
// %d data format at "01/02"
// %L log level e.g. INFO WARN ERROR
// %f function name and line number e.g. model.Get:121
// %i instance id
// %e deploy env e.g. dev uat fat prod
// %z zone
// %S full file name and line number: /a/b/c/d.go:23
// %s final file name element and line number: d.go:23
// %M log message and additional fields: key=value this is log message
func (h *StdoutHandler) SetFormat(format string) {
	h.render = newPatternRender(format)
}
