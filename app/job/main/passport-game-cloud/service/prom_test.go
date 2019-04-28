package service

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewInterval(t *testing.T) {
	name := "counter_1"
	NewInterval(&IntervalConfig{
		Name: name,
		Rate: 1000,
	})
	defer func() {
		expectErrStr := fmt.Sprintf("%s already exists", name)
		errStr := recover()
		if errStr != expectErrStr {
			t.Errorf("must painc when name is not unique, expcted %s but got %v", expectErrStr, errStr)
			t.FailNow()
		}
	}()
	NewInterval(&IntervalConfig{
		Name: name,
		Rate: 1000,
	})
}

func TestInterval_Prom(t *testing.T) {
	name := "counter_1"
	itv := NewInterval(&IntervalConfig{
		Name: name,
		Rate: 1000,
	})
	count := 10
	mtStr := time.Now().Format(_timeFormat)
	for i := 0; i < count; i++ {
		itv.Prom(context.TODO(), itv.MTS(context.TODO(), mtStr))
	}
	if itv.counter != int64(count) {
		t.Errorf("interval counter of %s should be %d after get mts for %d times, but now it is %d", name, count, count, itv.counter)
		t.FailNow()
	}
}
