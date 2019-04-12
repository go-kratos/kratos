package metric

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRollingGaugeAdd(t *testing.T) {
	size := 3
	bucketDuration := time.Second
	opts := RollingGaugeOpts{
		Size:           size,
		BucketDuration: bucketDuration,
	}
	r := NewRollingGauge(opts)
	listBuckets := func() [][]float64 {
		buckets := make([][]float64, 0)
		r.Reduce(func(i Iterator) float64 {
			for i.Next() {
				bucket := i.Bucket()
				buckets = append(buckets, bucket.Points)
			}
			return 0.0
		})
		return buckets
	}
	assert.Equal(t, [][]float64{{}, {}, {}}, listBuckets())
	r.Add(1)
	assert.Equal(t, [][]float64{{}, {}, {1}}, listBuckets())
	time.Sleep(time.Second)
	r.Add(2)
	r.Add(3)
	assert.Equal(t, [][]float64{{}, {1}, {2, 3}}, listBuckets())
	time.Sleep(time.Second)
	r.Add(4)
	r.Add(5)
	r.Add(6)
	assert.Equal(t, [][]float64{{1}, {2, 3}, {4, 5, 6}}, listBuckets())
	time.Sleep(time.Second)
	r.Add(7)
	assert.Equal(t, [][]float64{{2, 3}, {4, 5, 6}, {7}}, listBuckets())
}

func TestRollingGaugeReset(t *testing.T) {
	size := 3
	bucketDuration := time.Second
	opts := RollingGaugeOpts{
		Size:           size,
		BucketDuration: bucketDuration,
	}
	r := NewRollingGauge(opts)
	listBuckets := func() [][]float64 {
		buckets := make([][]float64, 0)
		r.Reduce(func(i Iterator) float64 {
			for i.Next() {
				bucket := i.Bucket()
				buckets = append(buckets, bucket.Points)
			}
			return 0.0
		})
		return buckets
	}
	r.Add(1)
	time.Sleep(time.Second)
	assert.Equal(t, [][]float64{{}, {1}}, listBuckets())
	time.Sleep(time.Second)
	assert.Equal(t, [][]float64{{1}}, listBuckets())
	time.Sleep(time.Second)
	assert.Equal(t, [][]float64{}, listBuckets())

	// cross window
	r.Add(1)
	time.Sleep(time.Second * 5)
	assert.Equal(t, [][]float64{}, listBuckets())
}

func TestRollingGaugeReduce(t *testing.T) {
	size := 3
	bucketDuration := time.Second
	opts := RollingGaugeOpts{
		Size:           size,
		BucketDuration: bucketDuration,
	}
	r := NewRollingGauge(opts)
	for x := 0; x < size; x = x + 1 {
		for i := 0; i <= x; i++ {
			r.Add(int64(i))
		}
		if x < size-1 {
			time.Sleep(bucketDuration)
		}
	}
	var result = r.Reduce(func(i Iterator) float64 {
		var result float64
		for i.Next() {
			bucket := i.Bucket()
			for _, point := range bucket.Points {
				result += point
			}
		}
		return result
	})
	if result != 4.0 {
		t.Fatalf("Validate sum of points. result: %f", result)
	}
}

func TestRollingGaugeDataRace(t *testing.T) {
	size := 3
	bucketDuration := time.Second
	opts := RollingGaugeOpts{
		Size:           size,
		BucketDuration: bucketDuration,
	}
	r := NewRollingGauge(opts)
	var stop = make(chan bool)
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				r.Add(rand.Int63())
				time.Sleep(time.Millisecond * 5)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				_ = r.Reduce(func(i Iterator) float64 {
					for i.Next() {
						bucket := i.Bucket()
						for range bucket.Points {
							continue
						}
					}
					return 0
				})
			}
		}
	}()
	time.Sleep(time.Second * 3)
	close(stop)
}

func BenchmarkRollingGaugeIncr(b *testing.B) {
	size := 10
	bucketDuration := time.Second
	opts := RollingGaugeOpts{
		Size:           size,
		BucketDuration: bucketDuration,
	}
	r := NewRollingGauge(opts)
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		r.Add(1.0)
	}
}

func BenchmarkRollingGaugeReduce(b *testing.B) {
	size := 10
	bucketDuration := time.Second
	opts := RollingGaugeOpts{
		Size:           size,
		BucketDuration: bucketDuration,
	}
	r := NewRollingGauge(opts)
	for i := 0; i <= 10; i++ {
		r.Add(1.0)
		time.Sleep(time.Millisecond * 500)
	}
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		var _ = r.Reduce(func(i Iterator) float64 {
			var result float64
			for i.Next() {
				bucket := i.Bucket()
				if len(bucket.Points) != 0 {
					result += bucket.Points[0]
				}
			}
			return result
		})
	}
}
