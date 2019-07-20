package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHistogramVec(t *testing.T) {
	histogramVec := NewHistogramVec(&HistogramVecOpts{
		Namespace: "test",
		Subsystem: "test",
		Name:      "test",
		Help:      "this is test metrics.",
		Labels:    []string{"name", "addr"},
		Buckets:   _defaultBuckets,
	})
	histogramVec.Observe(int64(1), "name1", "127.0.0.1")
	assert.Panics(t, func() {
		NewHistogramVec(&HistogramVecOpts{
			Namespace: "test",
			Subsystem: "test",
			Name:      "test",
			Help:      "this is test metrics.",
			Labels:    []string{"name", "addr"},
			Buckets:   _defaultBuckets,
		})
	}, "Expected to panic.")
	assert.NotPanics(t, func() {
		NewHistogramVec(&HistogramVecOpts{
			Namespace: "test",
			Subsystem: "test",
			Name:      "test2",
			Help:      "this is test metrics.",
			Labels:    []string{"name", "addr"},
			Buckets:   _defaultBuckets,
		})
	}, "Expected normal. no panic.")
}
