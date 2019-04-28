package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaohbaseMd5Key(t *testing.T) {
	convey.Convey("hbaseMd5Key", t, func(ctx convey.C) {
		var (
			aid = int64(1000)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := hbaseMd5Key(aid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoBaseUpStat(t *testing.T) {
	convey.Convey("BaseUpStat", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(1000)
			date = "2018-06-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.BaseUpStat(c, mid, date)
			ctx.Convey("Then err should be nil.stat should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
