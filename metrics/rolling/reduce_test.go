package rolling

import (
	"testing"

	"github.com/go-kratos/kratos/v2/metrics/point"
	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	opts := point.PointGaugeOpts{Size: 10}
	pointGauge := point.NewPointGauge(opts)
	for i := 0; i < opts.Size; i++ {
		pointGauge.Add(int64(i))
	}
	result := pointGauge.Reduce(Count)
	assert.Equal(t, float64(10), result, "validate count of pointGauge")
}
