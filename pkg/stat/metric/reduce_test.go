package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	opts := PointGaugeOpts{Size: 10}
	pointGauge := NewPointGauge(opts)
	for i := 0; i < opts.Size; i++ {
		pointGauge.Add(int64(i))
	}
	result := pointGauge.Reduce(Count)
	assert.Equal(t, float64(10), result, "validate count of pointGauge")
}
