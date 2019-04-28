package like

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikekeyRedDot(t *testing.T) {
	convey.Convey("keyRedDot", t, func(ctx convey.C) {
		var (
			mid = int64(55)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := keyRedDot(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeCacheRedDotTs(t *testing.T) {
	convey.Convey("CacheRedDotTs", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(77)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ts, err := d.CacheRedDotTs(c, mid)
			ctx.Convey("Then err should be nil.ts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ts, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeAddCacheRedDotTs(t *testing.T) {
	convey.Convey("AddCacheRedDotTs", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(77)
			ts  = int64(14585874)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddCacheRedDotTs(c, mid, ts)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
