package generator

import (
	"io"
)

// ErrorTracker tracks errors to the underlying writer, so that you can ignore
// them until you're ready to return.
type ErrorTracker struct {
	io.Writer
	err error
}

// NewErrorTracker makes a new error tracker; note that it implements io.Writer.
func NewErrorTracker(w io.Writer) *ErrorTracker {
	return &ErrorTracker{Writer: w}
}

// Write intercepts calls to Write.
func (et *ErrorTracker) Write(p []byte) (n int, err error) {
	if et.err != nil {
		return 0, et.err
	}
	n, err = et.Writer.Write(p)
	if err != nil {
		et.err = err
	}
	return n, err
}

// Error returns nil if no error has occurred, otherwise it returns the error.
func (et *ErrorTracker) Error() error {
	return et.err
}
