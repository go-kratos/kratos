package bnj

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestBnjAddCacheTimeFinish(t *testing.T) {
	convey.Convey("AddCacheTimeFinish", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			val = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddCacheTimeFinish(c, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

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

func TestBnjCacheLessTime(t *testing.T) {
	convey.Convey("CacheLessTime", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.CacheLessTime(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBnjAddCacheLessTime(t *testing.T) {
	convey.Convey("AddCacheLessTime", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			val = int64(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.AddCacheLessTime(c, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
