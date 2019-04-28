package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetUpTagIncome(t *testing.T) {
	convey.Convey("GetUpTagIncome", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = "2018-01-01"
			tagID = int64(11111)
			query = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_tag_income(tag_id, mid, av_id, date) VALUES(11111, 10, 10, '2018-01-01')")
			infos, err := d.GetUpTagIncome(c, date, tagID, query)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})

	convey.Convey("GetUpTagIncome query error", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			date  = "2018-01-01"
			tagID = int64(11111)
			query = "AND name = 'test'"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO up_tag_income(tag_id, mid, av_id, date) VALUES(11111, 10, 10, '2018-01-01')")
			_, err := d.GetUpTagIncome(c, date, tagID, query)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
