package dao

import (
	"context"
	"go-common/app/service/live/gift/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetWeekBagStatus(t *testing.T) {
	convey.Convey("GetWeekBagStatus", t, func(c convey.C) {
		var (
			ctx   = context.Background()
			uid   = int64(0)
			week  = int(0)
			level = int64(0)
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			res, err := d.GetWeekBagStatus(ctx, uid, week, level)
			c.Convey("Then err should be nil.res should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddWeekBag(t *testing.T) {
	convey.Convey("AddWeekBag", t, func(c convey.C) {
		var (
			ctx      = context.Background()
			uid      = int64(9527)
			week     = int(1)
			level    = int64(1)
			weekInfo = &model.BagGiftStatus{
				Status: 1,
				Gift: []*model.GiftInfo{
					{GiftID: 1, GiftNum: 2, ExpireAt: "7天"}, {GiftID: 2, GiftNum: 3, ExpireAt: "7天"},
				},
			}
		)
		c.Convey("When everything gose positive", func(c convey.C) {
			affected, err := d.AddWeekBag(ctx, uid, week, level, weekInfo)
			c.Convey("Then err should be nil.affected should not be nil.", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}
