package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoArchive(t *testing.T) {
	convey.Convey("Archive", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			aid = int64(10110788)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			a, err := d.Archive(c, aid)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoArchives(t *testing.T) {
	convey.Convey("Archives", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aids = []int64{10110788}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			a, err := d.Archives(c, aids)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoArticleMetas(t *testing.T) {
	convey.Convey("ArticleMetas", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aids = []int64{10110788}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.ArticleMetas(c, aids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoStats(t *testing.T) {
	convey.Convey("Stats", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			aids = []int64{10110788}
			ip   = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			a, err := d.Stats(c, aids, ip)
			ctx.Convey("Then err should be nil.a should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(a, convey.ShouldNotBeNil)
			})
		})
	})
}
