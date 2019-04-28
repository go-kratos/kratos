package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/up-rating/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetUpRatingCache(t *testing.T) {
	convey.Convey("GetUpRatingCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.GetUpRatingCache(c, mid)
			ctx.Convey("Then err should be nil.rating should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetUpRatingCache(t *testing.T) {
	convey.Convey("SetUpRatingCache", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mid    = int64(0)
			rating = &model.Rating{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetUpRatingCache(c, mid, rating)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoupRatingKey(t *testing.T) {
	convey.Convey("upRatingKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := upRatingKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
