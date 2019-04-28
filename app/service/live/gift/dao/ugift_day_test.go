package dao

import (
	"context"
	"go-common/app/service/live/gift/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetDayBagStatus(t *testing.T) {
	convey.Convey("GetDayBagStatus", t, func(c convey.C) {
		var (
			ctx  = context.Background()
			uid  = int64(0)
			date = ""
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			res, err := d.GetDayBagStatus(ctx, uid, date)
			c.Convey("Then err should be nil.res should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddDayBag(t *testing.T) {
	convey.Convey("AddDayBag", t, func(c convey.C) {
		var (
			ctx     = context.Background()
			uid     = int64(9527)
			date    = "2018-07-04 00:00:00"
			dayInfo = &model.BagGiftStatus{
				Status: 1,
				Gift: []*model.GiftInfo{
					{GiftID: 1, GiftNum: 2, ExpireAt: "今天"}, {GiftID: 2, GiftNum: 3, ExpireAt: "今天"},
				},
			}
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			affected, err := d.AddDayBag(ctx, uid, date, dayInfo)
			c.Convey("Then err should be nil.affected should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}
