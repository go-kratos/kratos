package log

import "testing"

func TestValue(t *testing.T) {
	logger := With(DefaultLogger, "caller", Caller(4))
	logger.Print("message", "helloworld")
}
