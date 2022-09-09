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
	l = With(l)
	l = With(l, "unpaired_key")
	c, ok := l.(*logger)
	assert.True(t, ok)
	assert.True(t, len(c.prefix) == 2)

	l = With(l, "ts", DefaultTimestamp, "caller", DefaultCaller, "unpaired_key")
	l = With(l, "ts", Timestamp(time.ANSIC), "caller", Caller(-1), "unpaired_key")
	c, ok = l.(*logger)
	assert.True(t, ok)
	assert.True(t, len(c.prefix) == 14)

	_ = l.Log(LevelInfo, "key1", "value1")
}

func TestWithContext(t *testing.T) {
	WithContext(context.Background(), nil)
}

func TestWithReplace(t *testing.T) {
	l := DefaultLogger
	l = WithReplace(l)
	l = WithReplace(l)
	l = WithReplace(l, "unpaired_key")
	c, ok := l.(*logger)
	assert.True(t, ok)
	assert.True(t, len(c.prefix) == 2)

	l = WithReplace(l, "ts", DefaultTimestamp, "caller", DefaultCaller, "unpaired_key")
	l = WithReplace(l, "ts", Timestamp(time.ANSIC), "caller", Caller(-1), "unpaired_key")
	c, ok = l.(*logger)
	assert.True(t, ok)
	assert.True(t, len(c.prefix) == 6)

	_ = l.Log(LevelInfo, "key1", "value1")
}
