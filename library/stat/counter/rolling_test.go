package counter

import (
	"testing"
	"time"
)

func TestRollingCounterMinInterval(t *testing.T) {
	count := NewRolling(time.Second/2, 10)
	tk1 := time.NewTicker(5 * time.Millisecond)
	defer tk1.Stop()
	for i := 0; i < 100; i++ {
		<-tk1.C
		count.Add(1)
	}

	v := count.Value()
	t.Logf("count value: %d", v)
	// 10% of error when bucket is 10
	if v < 90 || v > 110 {
		t.Errorf("expect value in [90-110] get %d", v)
	}
}
