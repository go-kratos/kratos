package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointGaugeAdd(t *testing.T) {
	opts := PointGaugeOpts{Size: 3}
	pointGauge := NewPointGauge(opts)
	listBuckets := func() [][]float64 {
		buckets := make([][]float64, 0)
		pointGauge.Reduce(func(i Iterator) float64 {
			for i.Next() {
				bucket := i.Bucket()
				buckets = append(buckets, bucket.Points)
			}
			return 0.0
		})
		return buckets
	}
	assert.Equal(t, [][]float64{{}, {}, {}}, listBuckets(), "Empty Buckets")
	pointGauge.Add(1)
	assert.Equal(t, [][]float64{{}, {}, {1}}, listBuckets(), "Point 1")
	pointGauge.Add(2)
	assert.Equal(t, [][]float64{{}, {1}, {2}}, listBuckets(), "Point 1, 2")
	pointGauge.Add(3)
	assert.Equal(t, [][]float64{{1}, {2}, {3}}, listBuckets(), "Point 1, 2, 3")
	pointGauge.Add(4)
	assert.Equal(t, [][]float64{{2}, {3}, {4}}, listBuckets(), "Point 2, 3, 4")
	pointGauge.Add(5)
	assert.Equal(t, [][]float64{{3}, {4}, {5}}, listBuckets(), "Point 3, 4, 5")
}

func TestPointGaugeReduce(t *testing.T) {
	opts := PointGaugeOpts{Size: 10}
	pointGauge := NewPointGauge(opts)
	for i := 0; i < opts.Size; i++ {
		pointGauge.Add(int64(i))
	}
	var _ = pointGauge.Reduce(func(i Iterator) float64 {
		idx := 0
		for i.Next() {
			bucket := i.Bucket()
			assert.Equal(t, bucket.Points[0], float64(idx), "validate points of pointGauge")
			idx++
		}
		return 0.0
	})
	assert.Equal(t, float64(9), pointGauge.Max(), "validate max of pointGauge")
	assert.Equal(t, float64(4.5), pointGauge.Avg(), "validate avg of pointGauge")
	assert.Equal(t, float64(0), pointGauge.Min(), "validate min of pointGauge")
	assert.Equal(t, float64(45), pointGauge.Sum(), "validate sum of pointGauge")
}
