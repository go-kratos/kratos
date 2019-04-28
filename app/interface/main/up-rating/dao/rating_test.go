package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoTaskStatus(t *testing.T) {
	convey.Convey("TaskStatus", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = "2018-11-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			status, err := d.TaskStatus(c, date)
			ctx.Convey("Then err should be nil.status should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(status, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpScore(t *testing.T) {
	convey.Convey("UpScore", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mon  = int(11)
			mid  = int64(1)
			date = "2018-11-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.UpScore(c, mon, mid, date)
			ctx.Convey("Then err should be nil.score should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
