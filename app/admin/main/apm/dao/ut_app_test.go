package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetAppCovCache(t *testing.T) {
	convey.Convey("SetAppCovCache", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetAppCovCache(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				// coverage, _ := d.GetAppCovCache(c, "main")
				// t.Logf("\nthe coverage is %.2f\n", coverage)
			})
		})
	})
}

func TestDaoGetAppCovCache(t *testing.T) {
	convey.Convey("GetAppCovCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			path = "main"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			coverage, err := d.GetAppCovCache(c, path)
			ctx.Convey("Then err should be nil.coverage should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(coverage, convey.ShouldNotBeNil)
				t.Logf("\nthe coverage is %.2f\n", coverage)
			})
		})
	})
}
