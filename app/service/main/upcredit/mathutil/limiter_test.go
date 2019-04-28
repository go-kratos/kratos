package mathutil

import (
	. "github.com/smartystreets/goconvey/convey"
	"math"
	"testing"
	"time"
)

func Test_limiter(t *testing.T) {
	Convey("limit interval", t, func() {
		var rate = 100.0
		var limit = NewLimiter(rate)
		var interval = 1.0 / rate
		var last time.Time

		for i := 0; i < 100; i++ {
			var t = <-limit.Token()
			if !last.IsZero() {
				var diff = t.Sub(last)
				So(math.Abs(diff.Seconds()-interval), ShouldBeLessThanOrEqualTo, 0.002)
			}
			last = t
		}
	})

	Convey("limit count", t, func() {
		var rate = 100.0
		var seconds = 10.0
		var limit = NewLimiter(rate)
		var expect = rate * seconds
		var timer = time.NewTimer(time.Duration(float64(time.Second) * seconds))
		var total = 0
		var run = true
		for run {
			select {
			case <-timer.C:
				run = false
			default:
				<-limit.Token()
				total++
			}
		}
		So(math.Abs(float64(total)-expect), ShouldBeLessThanOrEqualTo, rate*0.01)
	})
}
