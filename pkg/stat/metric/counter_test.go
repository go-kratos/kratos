package metric

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	counter := NewCounter(CounterOpts{})
	count := rand.Intn(100)
	for i := 0; i < count; i++ {
		counter.Add(1)
	}
	val := counter.Value()
	assert.Equal(t, val, int64(count))
}
