package retry

import (
	"math"
	"testing"
	"time"
)

func TestConstantBackoff(t *testing.T) {
	initialInterval := 100 * time.Millisecond
	for _, j := range []int{-1, 0, 1, 50} {
		maximumJitterInterval := time.Duration(j) * time.Millisecond
		ej := j
		if j < 0 {
			ej = 0
		}
		expectMaximumJitterInterval := time.Duration(ej) * time.Millisecond
		constantBackoff := NewConstantBackoff(initialInterval, maximumJitterInterval)
		for i := 0; i < 10000; i++ {
			bo := constantBackoff.Next(uint(i))
			if !(initialInterval <= bo && bo <= initialInterval+expectMaximumJitterInterval) {
				t.Errorf("initialInterval:%s,maximumJitterInterval:%s,got:%s\n", initialInterval, maximumJitterInterval, bo)
			}
		}
	}
}

func TestExponentialBackoff(t *testing.T) {
	initialInterval := 100 * time.Millisecond
	maxInterval := 1000 * time.Millisecond
	for _, j := range []int{-1, 0, 1, 50} {
		maximumJitterInterval := time.Duration(j) * time.Millisecond
		ej := j
		if j < 0 {
			ej = 0
		}
		expectMaximumJitterInterval := time.Duration(ej) * time.Millisecond
		exponentialBackoff := NewExponentialBackoff(2.0, initialInterval, maxInterval, maximumJitterInterval)
		for i := 0; i < 10000; i++ {
			bo := exponentialBackoff.Next(uint(i))
			if !(time.Duration(math.Min(float64(initialInterval/time.Millisecond)*math.Pow(2.0, float64(i)),
				float64(maxInterval/time.Millisecond)))*time.Millisecond <= bo && bo <= maxInterval+expectMaximumJitterInterval) {
				t.Errorf("exponentFactor:2.0,backoffInterval:%s,maxInterval:%s,maximumJitterInterval:%s,got:%s\n", initialInterval, maxInterval, maximumJitterInterval, bo)
			}
		}
	}
}

func TestRetrierFunc(t *testing.T) {
	linearRetrier := NewRetrierFunc(func(retry uint) time.Duration {
		if retry <= 0 {
			return 0 * time.Millisecond
		}
		return time.Duration(retry) * time.Millisecond
	})
	for i := 0; i < 10000; i++ {
		if linearRetrier.NextInterval(uint(i)) != time.Duration(i)*time.Millisecond {
			t.Errorf("excepted:%s,got%s", time.Duration(i)*time.Millisecond, linearRetrier.NextInterval(uint(i)))
		}
	}
}

func TestRetrier(t *testing.T) {
	initialInterval := 20 * time.Millisecond
	maxInterval := 100 * time.Millisecond
	maximumJitterInterval := time.Duration(5) * time.Millisecond

	constantBackoff := NewConstantBackoff(initialInterval, maximumJitterInterval)
	constantRetrier := NewRetrier(constantBackoff)
	bo := constantRetrier.NextInterval(1)
	if !(initialInterval <= bo) {
		t.Errorf("initialInterval:%s,maximumJitterInterval:%s,got:%s\n", initialInterval, maximumJitterInterval, bo)
	}
	exponentialBackoff := NewExponentialBackoff(2.0, initialInterval, maxInterval, maximumJitterInterval)
	exponentialRetrier := NewRetrier(exponentialBackoff)
	bo = exponentialRetrier.NextInterval(1)
	if !(40*time.Millisecond <= bo && bo <= maxInterval+maximumJitterInterval) {
		t.Errorf("exponentFactor:2.0,backoffInterval:%s,maxInterval:%s,maximumJitterInterval:%s,got:%s\n", initialInterval, maxInterval, maximumJitterInterval, bo)
	}
}

func TestNoRetrier(t *testing.T) {
	noRetrier := NewNoRetrier()
	if noRetrier.NextInterval(1) != time.Duration(0) {
		t.Errorf("excepted:%s,got:%s", time.Duration(0), noRetrier.NextInterval(1))
	}
}

func TestStrategy(t *testing.T) {
	strage := NewStrategy(3, NewNoRetrier(), []Condition{NewByCode(500, 504)})
	if strage.JudgeConditions(Resp{Code: 400}) {
		t.Errorf("expected:false,got:true")
	}
	if !strage.JudgeConditions(Resp{Code: 500}) {
		t.Errorf("expected:true,got:false")
	}
}
