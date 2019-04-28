package summary

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSummaryMinInterval(t *testing.T) {
	count := New(time.Second/2, 10)
	tk1 := time.NewTicker(5 * time.Millisecond)
	defer tk1.Stop()
	for i := 0; i < 100; i++ {
		<-tk1.C
		count.Add(2)
	}

	v, c := count.Value()
	t.Logf("count value: %d, %d\n", v, c)
	// 10% of error when bucket is 10
	if v < 190 || v > 210 {
		t.Errorf("expect value in [90-110] get %d", v)
	}
	// 10% of error when bucket is 10
	if c < 90 || c > 110 {
		t.Errorf("expect value in [90-110] get %d", v)
	}
}

func TestSummary(t *testing.T) {
	s := New(time.Second, 10)

	t.Run("add", func(t *testing.T) {
		s.Add(1)
		v, c := s.Value()
		assert.Equal(t, v, int64(1))
		assert.Equal(t, c, int64(1))
	})
	time.Sleep(time.Millisecond * 110)
	t.Run("add2", func(t *testing.T) {
		s.Add(1)
		v, c := s.Value()
		assert.Equal(t, v, int64(2))
		assert.Equal(t, c, int64(2))
	})
	time.Sleep(time.Millisecond * 900) // expire one bucket, 110 + 900
	t.Run("expire", func(t *testing.T) {
		v, c := s.Value()
		assert.Equal(t, v, int64(1))
		assert.Equal(t, c, int64(1))
		s.Add(1)
		v, c = s.Value()
		assert.Equal(t, v, int64(2)) // expire one bucket
		assert.Equal(t, c, int64(2)) // expire one bucket
	})
	time.Sleep(time.Millisecond * 1100)
	t.Run("expire_all", func(t *testing.T) {
		v, c := s.Value()
		assert.Equal(t, v, int64(0))
		assert.Equal(t, c, int64(0))
	})
	t.Run("reset", func(t *testing.T) {
		s.Reset()
		v, c := s.Value()
		assert.Equal(t, v, int64(0))
		assert.Equal(t, c, int64(0))
	})
}
