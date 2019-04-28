package dao

import (
	"context"
	resmdl "go-common/app/service/main/resource/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetIndexIconCache(t *testing.T) {
	convey.Convey("SetIndexIconCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = []*resmdl.IndexIcon{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetIndexIconCache(c, data)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoIndexIconCache(t *testing.T) {
	convey.Convey("IndexIconCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.IndexIconCache(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", res)
			})
		})
	})
}

func TestDaoIndexIconBakCache(t *testing.T) {
	convey.Convey("IndexIconBakCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.IndexIconBakCache(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.Printf("%+v", res)
			})
		})
	})
}
