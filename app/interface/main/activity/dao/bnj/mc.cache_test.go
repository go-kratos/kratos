package bnj

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestBnjCacheTimeFinish(t *testing.T) {
	convey.Convey("CacheTimeFinish", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CacheTimeFinish(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBnjDelCacheTimeFinish(t *testing.T) {
	convey.Convey("DelCacheTimeFinish", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheTimeFinish(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBnjDelCacheTimeLess(t *testing.T) {
	convey.Convey("DelCacheTimeLess", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.DelCacheTimeLess(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
