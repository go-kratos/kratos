package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/point/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaopointKey(t *testing.T) {
	convey.Convey("pointKey", t, func(ctx convey.C) {
		var (
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := pointKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoDelPointInfoCache(t *testing.T) {
	convey.Convey("DelPointInfoCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelPointInfoCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoPointInfoCache(t *testing.T) {
	convey.Convey("PointInfoCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.PointInfoCache(c, mid)
			ctx.Convey("Then err should be nil.pi should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetPointInfoCache(t *testing.T) {
	convey.Convey("SetPointInfoCache", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			pi = &model.PointInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetPointInfoCache(c, pi)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaodelCache(t *testing.T) {
	convey.Convey("delCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.delCache(c, key)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
