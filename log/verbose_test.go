package log

import "testing"

func TestVerbose(t *testing.T) {
	logger := With(DefaultLogger, "caller", DefaultCaller, "ts", DefaultTimestamp)
	v := NewVerbose(logger, 20)

	v.V(10).Log("foo", "bar1")
	v.V(20).Log("foo", "bar2")
	v.V(30).Log("foo", "bar3")
}
