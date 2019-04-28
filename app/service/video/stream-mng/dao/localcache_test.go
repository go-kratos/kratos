package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaogetLocalLiveStreamListKey(t *testing.T) {
	convey.Convey("getLocalLiveStreamListKey", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.getLocalLiveStreamListKey()
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLoadLiveStreamList(t *testing.T) {
	convey.Convey("LoadLiveStreamList", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rids = []int64{11891462}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := d.LoadLiveStreamList(c, rids)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoStoreLiveStreamList(t *testing.T) {
	convey.Convey("StoreLiveStreamList", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			d.StoreLiveStreamList()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
