package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestDao_CalPassScTime
func TestDao_CalPassScTime(t *testing.T) {
	Convey("CalPassScTime", t, func() {
		once.Do(startService)
		startTimes := map[int32]int32{0: 1530860100, 1: 1530860122}
		endTimes := map[int32]int32{0: 1533365702, 1: 1533970523}
		var tksPass []TicketPass
		tksPass = append(tksPass, TicketPass{Name: "test1", LinkScreens: []int32{0, 1}})
		tksPass = append(tksPass, TicketPass{Name: "test2", LinkScreens: []int32{0, 1}})
		res, err := d.CalPassScTime(startTimes, endTimes, tksPass)
		So(res, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})
}
