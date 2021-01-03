package file

import "testing"

func TestFormat(t *testing.T) {
	if ext := format("test.json"); ext != "json" {
		t.Errorf("no expected: %s", ext)
	}
}
