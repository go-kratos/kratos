package tag

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestTagGetArchiveIncome(t *testing.T) {
	convey.Convey("GetArchiveIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			query = ""
			limit = int(100)
			ctype = _video
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.db.Exec(c, "INSERT INTO av_income(av_id,date) VALUES(1,'2018-06-24') ON DUPLICATE KEY UPDATE date=VALUES(date)")
			archives, err := d.GetArchiveIncome(c, id, query, limit, ctype)
			ctx.Convey("Then err should be nil.archives should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(archives, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagGetUpIncome(t *testing.T) {
	convey.Convey("GetUpIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			query = ""
			limit = _video
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.GetUpIncome(c, id, query, limit)
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ups, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagGetAvIncomeStatis(t *testing.T) {
	convey.Convey("GetAvIncomeStatis", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			avs, err := d.GetAvIncomeStatis(c, id, limit)
			ctx.Convey("Then err should be nil.avs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(avs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagGetCmIncomeStatis(t *testing.T) {
	convey.Convey("GetCmIncomeStatis", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			id    = int64(0)
			limit = int(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cms, err := d.GetCmIncomeStatis(c, id, limit)
			ctx.Convey("Then err should be nil.cms should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(cms, convey.ShouldNotBeNil)
			})
		})
	})
}
