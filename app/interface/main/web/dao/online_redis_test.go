package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/web/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetOnlineListBakCache(t *testing.T) {
	convey.Convey("SetOnlineListBakCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = []*model.OnlineArc{{OnlineCount: 111}, {OnlineCount: 222}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetOnlineListBakCache(c, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoOnlineListBakCache(t *testing.T) {
	convey.Convey("OnlineListBakCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rs, err := d.OnlineListBakCache(c)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}
