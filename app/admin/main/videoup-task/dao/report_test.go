package dao

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTaskTooksByHalfHour(t *testing.T) {
	convey.Convey("TaskTooksByHalfHour", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			stime = time.Now()
			etime = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.TaskTooksByHalfHour(c, stime, etime)
			ctx.Convey("Then err should be nil.tooks should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
