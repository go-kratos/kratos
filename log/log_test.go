package log

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	l := DefaultLogger
	l = With(l)
	l = With(l, "error_key")
	c, ok := l.(*logger)
	assert.True(t, ok)
	assert.True(t, len(c.prefix) == 0)

	l = With(l, "ts", DefaultTimestamp, "caller", DefaultCaller, "error_key1")
	l = With(l, "ts", Timestamp(time.ANSIC), "caller", Caller(-1), "error_key2")
	c, ok = l.(*logger)
	assert.True(t, ok)
	assert.True(t, len(c.prefix) == 8)

	_ = l.Log(LevelInfo, "key1", "value1")
}

func TestWithContext(t *testing.T) {
	WithContext(context.Background(), nil)
}

func TestWithReplace(t *testing.T) {
	l := DefaultLogger
	l = WithReplace(l)
	l = WithReplace(l, "error_key")
	c, ok := l.(*logger)
	assert.True(t, ok)
	assert.True(t, len(c.prefix) == 0)

	l = WithReplace(l, "ts", DefaultTimestamp, "caller", DefaultCaller, "error_key1")
	l = WithReplace(l, "ts", Timestamp(time.ANSIC), "caller", Caller(-1), "error_key2")
	c, ok = l.(*logger)
	assert.True(t, ok)
	assert.True(t, len(c.prefix) == 4)

	_ = l.Log(LevelInfo, "key1", "value1")
}
