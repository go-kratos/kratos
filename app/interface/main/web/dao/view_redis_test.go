package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/web/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaokeyArchive(t *testing.T) {
	convey.Convey("keyArchive", t, func(ctx convey.C) {
		var (
			aid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyArchive(aid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetViewBakCache(t *testing.T) {
	convey.Convey("SetViewBakCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
			a   = &model.View{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetViewBakCache(c, aid, a)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoViewBakCache(t *testing.T) {
	convey.Convey("ViewBakCache", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rs, err := d.ViewBakCache(c, aid)
			ctx.Convey("Then err should be nil.rs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rs, convey.ShouldNotBeNil)
			})
		})
	})
}
