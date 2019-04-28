package breaker

import (
	"testing"
	"time"

	xtime "go-common/library/time"

	"github.com/stretchr/testify/assert"
)

func getSRE() Breaker {
	return NewGroup(&Config{
		Window:  xtime.Duration(1 * time.Second),
		Bucket:  10,
		Request: 100,
		K:       2,
	}).Get("")
}

func testSREClose(t *testing.T, b Breaker) {
	markSuccess(b, 80)
	assert.Equal(t, b.Allow(), nil)
	markSuccess(b, 120)
	assert.Equal(t, b.Allow(), nil)
}

func testSREOpen(t *testing.T, b Breaker) {
	markSuccess(b, 100)
	assert.Equal(t, b.Allow(), nil)
	markFailed(b, 10000000)
	assert.NotEqual(t, b.Allow(), nil)
}

func testSREHalfOpen(t *testing.T, b Breaker) {
	// failback
	assert.Equal(t, b.Allow(), nil)
	t.Run("allow single failed", func(t *testing.T) {
		markFailed(b, 10000000)
		assert.NotEqual(t, b.Allow(), nil)
	})
	time.Sleep(2 * time.Second)
	t.Run("allow single succeed", func(t *testing.T) {
		assert.Equal(t, b.Allow(), nil)
		markSuccess(b, 10000000)
		assert.Equal(t, b.Allow(), nil)
	})
}

func TestSRE(t *testing.T) {
	b := getSRE()
	testSREClose(t, b)

	b = getSRE()
	testSREOpen(t, b)

	b = getSRE()
	testSREHalfOpen(t, b)
}
