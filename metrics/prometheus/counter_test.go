package prometheus

import (
	"math/rand"
	"testing"

	"github.com/go-kratos/kratos/v2/metrics"
	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	counter := metrics.NewCounter(metrics.CounterOpts{})
	count := rand.Intn(100)
	for i := 0; i < count; i++ {
		counter.Add(1)
	}
	val := counter.Value()
	assert.Equal(t, val, int64(count))
}

func TestCounterVec(t *testing.T) {
	counterVec := NewCounterVec(&metrics.CounterVecOpts{
		Namespace: "test",
		Subsystem: "test",
		Name:      "test",
		Help:      "this is test metrics.",
		Labels:    []string{"name", "addr"},
	})
	counterVec.Inc("name1", "127.0.0.1")
	assert.Panics(t, func() {
		NewCounterVec(&metrics.CounterVecOpts{
			Namespace: "test",
			Subsystem: "test",
			Name:      "test",
			Help:      "this is test metrics.",
			Labels:    []string{"name", "addr"},
		})
	}, "Expected to panic.")
	assert.NotPanics(t, func() {
		NewCounterVec(&metrics.CounterVecOpts{
			Namespace: "test",
			Subsystem: "test",
			Name:      "test2",
			Help:      "this is test metrics.",
			Labels:    []string{"name", "addr"},
		})
	}, "Expected normal. no panic.")
}
