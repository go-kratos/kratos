// Package mocktrace this ut just make ci happay.
package mocktrace

import (
	"fmt"
	"testing"
)

func TestMockTrace(t *testing.T) {
	mocktrace := &MockTrace{}
	mocktrace.Inject(nil, nil, nil)
	mocktrace.Extract(nil, nil)

	root := mocktrace.New("test")
	root.Fork("", "")
	root.Follow("", "")
	root.Finish(nil)
	err := fmt.Errorf("test")
	root.Finish(&err)
	root.SetTag()
	root.SetLog()
	root.Visit(func(k, v string) {})
	root.SetTitle("")
}
