package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoListAvBreach(t *testing.T) {
	convey.Convey("ListAvBreach", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(1001)
			startTime = "2018-01-01"
			endTime   = "2019-01-01"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_breach_record(av_id, mid, cdate) VALUES(1000, 1001, '2018-06-01') ON DUPLICATE KEY UPDATE cdate = '2018-06-01'")
			records, err := d.ListAvBreach(c, mid, startTime, endTime)
			ctx.Convey("Then err should be nil.records should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(records, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAvBreachByType(t *testing.T) {
	convey.Convey("GetAvBreachByType", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(1002)
			begin = "2018-01-01"
			end   = "2019-01-01"
			typ   = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_breach_record(av_id, mid, cdate, ctype) VALUES(1001, 1002, '2018-06-01', 0) ON DUPLICATE KEY UPDATE cdate = '2018-06-01', ctype = 0")
			rs, err := d.GetAvBreachByType(c, mid, begin, end, typ)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetAvBreachs(t *testing.T) {
	convey.Convey("GetAvBreachs", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			avs   = []int64{1000}
			ctype = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO av_breach_record(av_id, mid, cdate, ctype) VALUES(1001, 1002, '2018-06-01', 0) ON DUPLICATE KEY UPDATE cdate = '2018-06-01', ctype = 0")
			breachs, err := d.GetAvBreachs(c, avs, ctype)
			ctx.Convey("Then err should be nil.breachs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(breachs, convey.ShouldNotBeNil)
			})
		})
	})
}
