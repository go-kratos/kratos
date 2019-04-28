package counter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGaugeCounter(t *testing.T) {
	key := "test"
	g := &Group{
		New: func() Counter {
			return NewGauge()
		},
	}
	g.Add(key, 1)
	g.Add(key, 2)
	g.Add(key, 3)
	g.Add(key, -1)
	assert.Equal(t, g.Value(key), int64(5))
	g.Reset(key)
	assert.Equal(t, g.Value(key), int64(0))
}

func TestRollingCounter(t *testing.T) {
	key := "test"
	g := &Group{
		New: func() Counter {
			return NewRolling(time.Second, 10)
		},
	}

	t.Run("add_key_b1", func(t *testing.T) {
		g.Add(key, 1)
		assert.Equal(t, g.Value(key), int64(1))
	})
	time.Sleep(time.Millisecond * 110)
	t.Run("add_key_b2", func(t *testing.T) {
		g.Add(key, 1)
		assert.Equal(t, g.Value(key), int64(2))
	})
	time.Sleep(time.Millisecond * 900) // expire one bucket, 110 + 900
	t.Run("expire_b1", func(t *testing.T) {
		assert.Equal(t, g.Value(key), int64(1))
		g.Add(key, 1)
		assert.Equal(t, g.Value(key), int64(2)) // expire one bucket
	})
	time.Sleep(time.Millisecond * 1100)
	t.Run("expire_all", func(t *testing.T) {
		assert.Equal(t, g.Value(key), int64(0))
	})
	t.Run("reset", func(t *testing.T) {
		g.Reset(key)
		assert.Equal(t, g.Value(key), int64(0))
	})
}
