package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGaugeAdd(t *testing.T) {
	gauge := NewGauge(GaugeOpts{})
	gauge.Add(100)
	gauge.Add(-50)
	val := gauge.Value()
	assert.Equal(t, val, int64(50))
}

func TestGaugeSet(t *testing.T) {
	gauge := NewGauge(GaugeOpts{})
	gauge.Add(100)
	gauge.Set(50)
	val := gauge.Value()
	assert.Equal(t, val, int64(50))
}
