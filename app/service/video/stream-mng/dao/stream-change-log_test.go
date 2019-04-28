package dao

import (
	"context"
	"go-common/app/service/video/stream-mng/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInsertChangeLog(t *testing.T) {
	convey.Convey("InsertChangeLog", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			change = &model.StreamChangeLog{
				RoomID:      11891462,
				FromOrigin:  1,
				ToOrigin:    2,
				Source:      "app",
				OperateName: "yy",
				Reason:      "auto",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.InsertChangeLog(c, change)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGetChangeLogByRoomID(t *testing.T) {
	convey.Convey("GetChangeLogByRoomID", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			rid   = int64(11891462)
			limit = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			infos, err := d.GetChangeLogByRoomID(c, rid, limit)
			ctx.Convey("Then err should be nil.infos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(infos, convey.ShouldNotBeNil)
			})
		})
	})
}
