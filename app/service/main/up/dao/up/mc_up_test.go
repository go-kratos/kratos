package up

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpupCacheKey(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("upCacheKey", t, func(ctx convey.C) {
		p1 := upCacheKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestUpupSwitchKey(t *testing.T) {
	var (
		mid = int64(0)
	)
	convey.Convey("upSwitchKey", t, func(ctx convey.C) {
		p1 := upSwitchKey(mid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestUpDelCpUp(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(0)
	)
	convey.Convey("DelCpUp", t, func(ctx convey.C) {
		d.DelCpUp(c, mid)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}
