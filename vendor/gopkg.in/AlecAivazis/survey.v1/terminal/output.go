// +build !windows

package terminal

import (
	"io"
)

// Returns special stdout, which converts escape sequences to Windows API calls
// on Windows environment.
func NewAnsiStdout(out FileWriter) io.Writer {
	return out
}

// Returns special stderr, which converts escape sequences to Windows API calls
// on Windows environment.
func NewAnsiStderr(out FileWriter) io.Writer {
	return out
}
