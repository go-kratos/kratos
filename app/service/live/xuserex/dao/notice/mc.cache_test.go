package notice

import (
	"context"
	"go-common/app/service/live/xuserex/model/roomNotice"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

// TestRoomNoticeAddCacheMonthConsume test.
func TestRoomNoticeAddCacheMonthConsume(t *testing.T) {
	convey.Convey("AddCacheMonthConsume", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			id       = int64(10000)
			targetID = int64(1008)
			date     = "20190101"
			value    = &roomNotice.MonthConsume{
				Uid:    id,
				Ruid:   targetID,
				Amount: 10,
				Date:   -1,
			}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddCacheMonthConsume(c, id, targetID, date, value)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

// TestRoomNoticeCacheMonthConsume auto test.
func TestRoomNoticeCacheMonthConsume(t *testing.T) {
	convey.Convey("CacheMonthConsume", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			id       = int64(10000)
			targetID = int64(1008)
			date     = "20190101"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CacheMonthConsume(c, id, targetID, date)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
