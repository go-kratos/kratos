package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPushDatas(t *testing.T) {
	convey.Convey("PushDatas", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			curtime = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.PushDatas(c, curtime)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdatePushData(t *testing.T) {
	convey.Convey("UpdatePushData", t, func(ctx convey.C) {
		var (
			c              = context.Background()
			status         = int8(0)
			progressStatus = int8(0)
			pushedCount    = int32(0)
			errcode        = int64(0)
			data           = int64(0)
			id             = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdatePushData(c, status, progressStatus, pushedCount, errcode, data, id)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
