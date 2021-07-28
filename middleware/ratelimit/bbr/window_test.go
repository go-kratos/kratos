package bbr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWindowResetWindow(t *testing.T) {
	opts := WindowOpts{Size: 3}
	window := NewWindow(opts)
	for i := 0; i < opts.Size; i++ {
		window.Append(i, 1.0)
	}
	window.ResetWindow()
	for i := 0; i < opts.Size; i++ {
		assert.Equal(t, len(window.Bucket(i).Points), 0)
	}
}

func TestWindowResetBucket(t *testing.T) {
	opts := WindowOpts{Size: 3}
	window := NewWindow(opts)
	for i := 0; i < opts.Size; i++ {
		window.Append(i, 1.0)
	}
	window.ResetBucket(1)
	assert.Equal(t, len(window.Bucket(1).Points), 0)
	assert.Equal(t, window.Bucket(0).Points[0], float64(1.0))
	assert.Equal(t, window.Bucket(2).Points[0], float64(1.0))
}

func TestWindowResetBuckets(t *testing.T) {
	opts := WindowOpts{Size: 3}
	window := NewWindow(opts)
	for i := 0; i < opts.Size; i++ {
		window.Append(i, 1.0)
	}
	window.ResetBuckets(0, 3)
	for i := 0; i < opts.Size; i++ {
		assert.Equal(t, len(window.Bucket(i).Points), 0)
	}
}

func TestWindowAppend(t *testing.T) {
	opts := WindowOpts{Size: 3}
	window := NewWindow(opts)
	for i := 0; i < opts.Size; i++ {
		window.Append(i, 1.0)
	}
	for i := 0; i < opts.Size; i++ {
		assert.Equal(t, window.Bucket(i).Points[0], float64(1.0))
	}
}

func TestWindowAdd(t *testing.T) {
	opts := WindowOpts{Size: 3}
	window := NewWindow(opts)
	window.Append(0, 1.0)
	window.Add(0, 1.0)
	assert.Equal(t, window.Bucket(0).Points[0], float64(2.0))
}

func TestWindowSize(t *testing.T) {
	opts := WindowOpts{Size: 3}
	window := NewWindow(opts)
	assert.Equal(t, window.Size(), 3)
}
