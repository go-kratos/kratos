package upcrm

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmScoreQueryHistory(t *testing.T) {
	convey.Convey("ScoreQueryHistory", t, func(ctx convey.C) {
		var (
			scoreType = int(0)
			date      = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.ScoreQueryHistory(scoreType, date)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmGetLastHistory(t *testing.T) {
	convey.Convey("GetLastHistory", t, func(ctx convey.C) {
		var (
			scoreType = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			lastHistoryDate, err := d.GetLastHistory(scoreType)
			ctx.Convey("Then err should be nil.lastHistoryDate should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(lastHistoryDate, convey.ShouldNotBeNil)
			})
		})
	})
}
