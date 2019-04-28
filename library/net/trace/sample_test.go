package trace

import (
	"testing"
)

func TestProbabilitySampling(t *testing.T) {
	sampler := newSampler(0.001)
	t.Run("test one operationName", func(t *testing.T) {
		sampled, probability := sampler.IsSampled(0, "test123")
		if !sampled || probability != 1 {
			t.Errorf("expect sampled and probability == 1 get: %v %f", sampled, probability)
		}
	})
	t.Run("test probability", func(t *testing.T) {
		sampler.IsSampled(0, "test_opt_2")
		count := 0
		for i := 0; i < 100000; i++ {
			sampled, _ := sampler.IsSampled(0, "test_opt_2")
			if sampled {
				count++
			}
		}
		if count < 80 || count > 120 {
			t.Errorf("expect count between 80~120 get %d", count)
		}
	})
}

func BenchmarkProbabilitySampling(b *testing.B) {
	sampler := newSampler(0.001)
	for i := 0; i < b.N; i++ {
		sampler.IsSampled(0, "test_opt_xxx")
	}
}
