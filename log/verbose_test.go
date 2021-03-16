package log

import "testing"

func TestVerbose(t *testing.T) {
	v := NewVerbose(DefaultLogger, 20)

	v.V(10).Print("foo", "bar1")
	v.V(20).Print("foo", "bar2")
	v.V(30).Print("foo", "bar3")
}
