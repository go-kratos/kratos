package income

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIncomeGetAvByBgm(t *testing.T) {
	convey.Convey("GetAvByBgm", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			sid  = int64(1003)
			from = "2018-01-01"
			to   = "2019-01-03"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO bgm_income(sid,mid,date) VALUES(1003, 1000, '2018-05-01')")
			avs, err := d.GetAvByBgm(c, sid, from, to)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestIncomeGetBgmIncome(t *testing.T) {
	convey.Convey("GetBgmIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			query = ""
			from  = "2018-01-02"
			to    = "2018-01-03"
			limit = int(10)
			typ   = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			bgms, err := d.GetBgmIncome(c, id, query, from, to, limit, typ)
			ctx.Convey("Then err should be nil.bgms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(bgms, convey.ShouldNotBeNil)
			})
		})
	})
}
