package log

import (
	"context"

	pkgerr "github.com/pkg/errors"
)

const (
	_log = "log"
)

// Handler is used to handle log events, outputting them to
// stdio or sending them to remote services. See the "handlers"
// directory for implementations.
//
// It is left up to Handlers to implement thread-safety.
type Handler interface {
	// Log handle log
	// variadic D is k-v struct represent log content
	Log(context.Context, Level, ...D)

	// SetFormat set render format on log output
	// see StdoutHandler.SetFormat for detail
	SetFormat(string)

	// Close handler
	Close() error
}

// Handlers .
type Handlers []Handler

// Log handlers logging.
func (hs Handlers) Log(c context.Context, lv Level, d ...D) {
	for _, h := range hs {
		h.Log(c, lv, d...)
	}
}

// Close close resource.
func (hs Handlers) Close() (err error) {
	for _, h := range hs {
		if e := h.Close(); e != nil {
			err = pkgerr.WithStack(e)
		}
	}
	return
}

// SetFormat .
func (hs Handlers) SetFormat(format string) {
	for _, h := range hs {
		h.SetFormat(format)
	}
}
