package cache

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestCacheRefreshUpTypeAsync(t *testing.T) {
	convey.Convey("RefreshUpTypeAsync", t, func(ctx convey.C) {
		var (
			tm = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			RefreshUpTypeAsync(tm)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestCacheRefreshUpType(t *testing.T) {
	convey.Convey("RefreshUpType", t, func(ctx convey.C) {
		var (
			tm = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			RefreshUpType(tm)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestCacheGetTidName(t *testing.T) {
	convey.Convey("GetTidName", t, func(ctx convey.C) {
		var (
			tid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := GetTidName(tid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
